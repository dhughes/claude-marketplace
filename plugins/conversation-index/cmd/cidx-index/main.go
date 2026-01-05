package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/doughughes/claude-marketplace/plugins/conversation-index/internal/db"
	"github.com/doughughes/claude-marketplace/plugins/conversation-index/internal/indexer"
	"github.com/doughughes/claude-marketplace/plugins/conversation-index/internal/shared"
)

func main() {
	// Parse command line flags
	fullReindex := flag.Bool("full-reindex", false, "Drop existing index and reindex all conversations")
	flag.BoolVar(fullReindex, "f", false, "Drop existing index and reindex all conversations (shorthand)")

	flag.Parse()

	// Open database
	database, err := db.Open(shared.DBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// Create indexer
	idx := indexer.NewIndexer(database, shared.ProjectsDir)

	// Run indexing
	if err := idx.IndexAll(*fullReindex); err != nil {
		fmt.Fprintf(os.Stderr, "Error indexing conversations: %v\n", err)
		os.Exit(1)
	}
}
