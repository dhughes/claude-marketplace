#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const initSqlJs = require('sql.js-fts5');
const os = require('os');

// Configuration
const CLAUDE_DIR = path.join(os.homedir(), '.claude');
const PROJECTS_DIR = path.join(CLAUDE_DIR, 'projects');
const DB_PATH = path.join(CLAUDE_DIR, 'conversation-index.db');

// Load database from file or create new one
async function loadDB() {
  // In Node.js, locateFile can be omitted - the library finds wasm automatically
  const SQL = await initSqlJs();

  if (fs.existsSync(DB_PATH)) {
    const buffer = fs.readFileSync(DB_PATH);
    return new SQL.Database(buffer);
  }

  return new SQL.Database();
}

// Save database to file
function saveDB(db) {
  const data = db.export();
  const buffer = Buffer.from(data);
  fs.writeFileSync(DB_PATH, buffer);
}

// Initialize database schema
function initDB(db) {
  db.exec(`
    CREATE TABLE IF NOT EXISTS conversations (
      uuid TEXT PRIMARY KEY,
      project_path TEXT NOT NULL,
      encoded_path TEXT NOT NULL,
      created_at TEXT NOT NULL,
      last_updated TEXT NOT NULL,
      message_count INTEGER DEFAULT 0
    );

    CREATE TABLE IF NOT EXISTS messages (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      conversation_uuid TEXT NOT NULL,
      timestamp TEXT NOT NULL,
      role TEXT NOT NULL,
      content TEXT,
      FOREIGN KEY (conversation_uuid) REFERENCES conversations(uuid)
    );

    CREATE VIRTUAL TABLE IF NOT EXISTS messages_fts USING fts5(
      conversation_uuid,
      content,
      content=messages,
      content_rowid=id
    );

    CREATE TABLE IF NOT EXISTS index_state (
      conversation_uuid TEXT PRIMARY KEY,
      last_indexed_line INTEGER DEFAULT 0,
      last_modified_time TEXT
    );

    -- Triggers to keep FTS in sync
    CREATE TRIGGER IF NOT EXISTS messages_ai AFTER INSERT ON messages BEGIN
      INSERT INTO messages_fts(rowid, conversation_uuid, content)
      VALUES (new.id, new.conversation_uuid, new.content);
    END;

    CREATE TRIGGER IF NOT EXISTS messages_ad AFTER DELETE ON messages BEGIN
      DELETE FROM messages_fts WHERE rowid = old.id;
    END;

    CREATE TRIGGER IF NOT EXISTS messages_au AFTER UPDATE ON messages BEGIN
      DELETE FROM messages_fts WHERE rowid = old.id;
      INSERT INTO messages_fts(rowid, conversation_uuid, content)
      VALUES (new.id, new.conversation_uuid, new.content);
    END;
  `);
}

// Decode project path (reverse the encoding: -Users-foo-bar -> /Users/foo/bar)
// Note: This is lossy if directory names contain dashes
function decodeProjectPath(encoded) {
  if (!encoded.startsWith('-')) return encoded;
  return '/' + encoded.slice(1).replace(/-/g, '/');
}

// Encode project path (forward encoding: /Users/foo/bar -> -Users-foo-bar)
function encodeProjectPath(decoded) {
  if (decoded.startsWith('/')) {
    return '-' + decoded.slice(1).replace(/\//g, '-');
  }
  return decoded.replace(/\//g, '-');
}

// Extract searchable content from a JSONL line
function extractContent(line) {
  try {
    const entry = JSON.parse(line);
    const results = [];

    // User messages
    if (entry.type === 'user' && entry.message?.content) {
      results.push({
        role: 'user',
        content: entry.message.content,
        timestamp: entry.timestamp
      });
    }

    // Assistant messages (may contain text and/or tool calls)
    if (entry.type === 'assistant' && entry.message?.content) {
      const content = entry.message.content;

      if (Array.isArray(content)) {
        for (const item of content) {
          // Text responses
          if (item.type === 'text' && item.text) {
            results.push({
              role: 'assistant',
              content: item.text,
              timestamp: entry.timestamp
            });
          }

          // Tool calls
          if (item.type === 'tool_use' && item.name) {
            const parts = [`Tool: ${item.name}`];
            if (item.input) {
              // Extract interesting fields from tool input
              if (item.input.file_path) parts.push(`File: ${item.input.file_path}`);
              if (item.input.pattern) parts.push(`Pattern: ${item.input.pattern}`);
              if (item.input.description) parts.push(item.input.description);
              if (item.input.prompt) parts.push(item.input.prompt);
              if (item.input.command) parts.push(`Command: ${item.input.command}`);
            }
            results.push({
              role: 'tool',
              content: parts.join(' '),
              timestamp: entry.timestamp
            });
          }
        }
      }
    }

    return results.length > 0 ? results : null;
  } catch (e) {
    return null;
  }
}

// Index a single conversation file
function indexConversation(db, conversationPath, projectPath, encodedPath) {
  const uuid = path.basename(conversationPath, '.jsonl');
  const stats = fs.statSync(conversationPath);
  const lastModified = stats.mtime.toISOString();

  // Check if we need to index
  const stmt = db.prepare('SELECT last_indexed_line, last_modified_time FROM index_state WHERE conversation_uuid = ?');
  stmt.bind([uuid]);
  const state = stmt.step() ? stmt.getAsObject() : null;
  stmt.free();

  if (state && state.last_modified_time === lastModified) {
    // No changes since last index
    return { indexed: 0, skipped: true };
  }

  const lines = fs.readFileSync(conversationPath, 'utf-8').split('\n').filter(l => l.trim());

  // Detect rollback: file has fewer lines than we last indexed
  if (state && lines.length < state.last_indexed_line) {
    // Rollback detected - clear this conversation's index and re-index from scratch
    db.run('DELETE FROM messages WHERE conversation_uuid = ?', [uuid]);
    db.run('DELETE FROM index_state WHERE conversation_uuid = ?', [uuid]);
    db.run('UPDATE conversations SET message_count = 0 WHERE uuid = ?', [uuid]);
    // Continue to re-index from line 0
  }

  const startLine = state && lines.length >= state.last_indexed_line ? state.last_indexed_line : 0;

  if (lines.length <= startLine) {
    // No new lines
    return { indexed: 0, skipped: true };
  }

  // Get or create conversation record
  const convStmt = db.prepare('SELECT * FROM conversations WHERE uuid = ?');
  convStmt.bind([uuid]);
  let conversation = convStmt.step() ? convStmt.getAsObject() : null;
  convStmt.free();

  // Parse first line to get creation timestamp and actual project path
  let createdAt = conversation?.created_at || new Date().toISOString();
  let actualProjectPath = projectPath; // fallback to decoded path

  if (!conversation) {
    try {
      const firstEntry = JSON.parse(lines[0]);
      createdAt = firstEntry.timestamp || createdAt;
      // Extract actual cwd from JSONL - this is the accurate project path
      actualProjectPath = firstEntry.cwd || projectPath;
    } catch (e) {}
  } else {
    // Use existing project path from database
    actualProjectPath = conversation.project_path;
  }

  // Use UPSERT to handle race conditions and rollbacks
  db.run(`
    INSERT INTO conversations (uuid, project_path, encoded_path, created_at, last_updated, message_count)
    VALUES (?, ?, ?, ?, ?, COALESCE((SELECT message_count FROM conversations WHERE uuid = ?), 0))
    ON CONFLICT(uuid) DO UPDATE SET
      project_path = excluded.project_path,
      encoded_path = excluded.encoded_path,
      last_updated = excluded.last_updated
  `, [uuid, actualProjectPath, encodedPath, createdAt, lastModified, uuid]);

  // Index new lines
  const insertMessage = db.prepare(`
    INSERT INTO messages (conversation_uuid, timestamp, role, content)
    VALUES (?, ?, ?, ?)
  `);

  let indexed = 0;
  for (let i = startLine; i < lines.length; i++) {
    const extracted = extractContent(lines[i]);
    if (extracted) {
      for (const item of extracted) {
        const timestamp = item.timestamp ? String(item.timestamp) : new Date().toISOString();
        const content = item.content ? String(item.content) : '';
        insertMessage.bind([uuid, timestamp, item.role, content]);
        insertMessage.step();
        insertMessage.reset();
        indexed++;
      }
    }
  }
  insertMessage.free();

  // Update index state
  db.run(`
    INSERT OR REPLACE INTO index_state (conversation_uuid, last_indexed_line, last_modified_time)
    VALUES (?, ?, ?)
  `, [uuid, lines.length, lastModified]);

  // Update message count
  db.run('UPDATE conversations SET message_count = message_count + ? WHERE uuid = ?', [indexed, uuid]);

  return { indexed, skipped: false };
}

// Main indexing function
async function indexAllConversations(fullReindex = false) {
  const db = await loadDB();
  initDB(db);

  if (fullReindex) {
    console.log('Performing full reindex...');
    db.exec('DELETE FROM messages; DELETE FROM conversations; DELETE FROM index_state;');
  }

  const startTime = Date.now();
  let totalIndexed = 0;
  let totalSkipped = 0;
  let totalConversations = 0;

  // Scan all project directories
  if (!fs.existsSync(PROJECTS_DIR)) {
    console.log('No projects directory found');
    db.close();
    return;
  }

  const projectDirs = fs.readdirSync(PROJECTS_DIR);

  for (const projectDir of projectDirs) {
    const projectPath = path.join(PROJECTS_DIR, projectDir);
    const decodedPath = decodeProjectPath(projectDir);

    if (!fs.statSync(projectPath).isDirectory()) continue;

    const conversations = fs.readdirSync(projectPath).filter(f => f.endsWith('.jsonl'));

    for (const conversation of conversations) {
      const conversationPath = path.join(projectPath, conversation);
      const result = indexConversation(db, conversationPath, decodedPath, projectDir);

      totalConversations++;
      if (result.skipped) {
        totalSkipped++;
      } else {
        totalIndexed += result.indexed;
      }
    }
  }

  // Save database to file
  saveDB(db);
  db.close();

  const elapsed = Date.now() - startTime;
  console.log(`Indexed ${totalIndexed} messages from ${totalConversations} conversations (${totalSkipped} skipped) in ${elapsed}ms`);
}

// CLI
const args = process.argv.slice(2);
const fullReindex = args.includes('--full-reindex') || args.includes('-f');

indexAllConversations(fullReindex).catch(err => {
  console.error('Error indexing conversations:', err);
  process.exit(1);
});
