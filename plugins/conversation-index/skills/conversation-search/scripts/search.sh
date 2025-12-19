#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const Database = require('better-sqlite3');
const os = require('os');

// Configuration
const CLAUDE_DIR = path.join(os.homedir(), '.claude');
const DB_PATH = path.join(CLAUDE_DIR, 'conversation-index.db');

// Encode project path (forward encoding: /Users/foo/bar -> -Users-foo-bar)
function encodeProjectPath(decoded) {
  if (decoded.startsWith('/')) {
    return '-' + decoded.slice(1).replace(/\//g, '-');
  }
  return decoded.replace(/\//g, '-');
}

function search(query, scope = 'current_project', currentProject = null, limit = 100) {
  if (!fs.existsSync(DB_PATH)) {
    return {
      error: 'Index not initialized. Run indexer first.',
      matches: []
    };
  }

  const db = new Database(DB_PATH, { readonly: true });

  let sql = `
    SELECT
      c.uuid,
      c.project_path,
      c.encoded_path,
      c.created_at,
      c.last_updated,
      c.message_count,
      GROUP_CONCAT(m.content, ' ') as content_preview,
      messages_fts.rank as relevance_score
    FROM messages_fts
    JOIN messages m ON messages_fts.rowid = m.id
    JOIN conversations c ON m.conversation_uuid = c.uuid
    WHERE messages_fts MATCH ?
  `;

  const params = [query];

  if (scope === 'current_project' && currentProject) {
    sql += ` AND c.encoded_path = ?`;
    params.push(encodeProjectPath(currentProject));
  }

  sql += `
    GROUP BY c.uuid
    ORDER BY relevance_score DESC
    LIMIT ?
  `;
  params.push(limit);

  const results = db.prepare(sql).all(...params);

  // Generate summaries by extracting first user message
  const summaries = results.map(result => {
    const firstMessage = db.prepare(`
      SELECT content FROM messages
      WHERE conversation_uuid = ? AND role = 'user'
      ORDER BY timestamp ASC
      LIMIT 1
    `).get(result.uuid);

    // Truncate content for preview
    let summary = firstMessage?.content || 'No summary available';
    if (summary.length > 150) {
      summary = summary.substring(0, 147) + '...';
    }

    return {
      uuid: result.uuid,
      project_path: result.project_path,
      encoded_path: result.encoded_path,
      created_at: result.created_at,
      last_updated: result.last_updated,
      message_count: result.message_count,
      summary,
      relevance_score: Math.abs(result.relevance_score)
    };
  });

  db.close();

  return {
    query,
    scope,
    current_project: currentProject,
    total_matches: summaries.length,
    matches: summaries
  };
}

// CLI
function main() {
  const args = process.argv.slice(2);

  if (args.length === 0 || args.includes('--help') || args.includes('-h')) {
    console.log(`Usage: search.sh [options] <query>

Options:
  --scope <current_project|all_projects>  Search scope (default: current_project)
  --project <path>                         Current project path for scoping
  --limit <number>                         Maximum results (default: 100)
  --json                                   Output as JSON

Examples:
  search.sh "authentication system"
  search.sh --scope all_projects "bug fix"
  search.sh --project "/Users/doug/code/app" "API"
`);
    process.exit(0);
  }

  let query = '';
  let scope = 'current_project';
  let currentProject = process.cwd();
  let limit = 100;
  let jsonOutput = false;

  for (let i = 0; i < args.length; i++) {
    if (args[i] === '--scope') {
      scope = args[++i];
    } else if (args[i] === '--project') {
      currentProject = args[++i];
    } else if (args[i] === '--limit') {
      limit = parseInt(args[++i]);
    } else if (args[i] === '--json') {
      jsonOutput = true;
    } else {
      query = args[i];
    }
  }

  if (!query) {
    console.error('Error: Query required');
    process.exit(1);
  }

  const results = search(query, scope, currentProject, limit);

  if (jsonOutput) {
    console.log(JSON.stringify(results, null, 2));
  } else {
    if (results.error) {
      console.error(results.error);
      process.exit(1);
    }

    console.log(`Found ${results.total_matches} conversation(s) matching "${results.query}"\n`);

    if (results.matches.length === 0) {
      console.log('No matches found.');
    } else {
      results.matches.forEach((match, idx) => {
        console.log(`${idx + 1}. UUID: ${match.uuid}`);
        console.log(`   Project: ${match.project_path}`);
        console.log(`   Created: ${new Date(match.created_at).toLocaleString()}`);
        console.log(`   Messages: ${match.message_count}`);
        console.log(`   Summary: ${match.summary}`);
        console.log(`   Relevance: ${match.relevance_score.toFixed(2)}`);
        console.log('');
      });
    }
  }
}

main();
