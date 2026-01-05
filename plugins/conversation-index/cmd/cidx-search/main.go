package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/doughughes/claude-marketplace/plugins/conversation-index/internal/db"
	"github.com/doughughes/claude-marketplace/plugins/conversation-index/internal/shared"
)

func main() {
	// Parse command line flags
	scope := flag.String("scope", "current_project", "Search scope: current_project or all_projects")
	project := flag.String("project", "", "Current project path (default: cwd)")
	limit := flag.Int("limit", 100, "Maximum results")
	jsonOutput := flag.Bool("json", false, "Output as JSON")
	help := flag.Bool("help", false, "Show help")
	flag.BoolVar(help, "h", false, "Show help (shorthand)")

	flag.Parse()

	// Show help
	if *help || flag.NArg() == 0 {
		fmt.Println(`Usage: search [options] <query>

Options:
  --scope <current_project|all_projects>  Search scope (default: current_project)
  --project <path>                         Current project path for scoping
  --limit <number>                         Maximum results (default: 100)
  --json                                   Output as JSON

Examples:
  search "authentication system"
  search --scope all_projects "bug fix"
  search --project "/Users/doug/code/app" "API"`)
		os.Exit(0)
	}

	query := flag.Arg(0)

	// Default project to current working directory
	if *project == "" {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
			os.Exit(1)
		}
		*project = shared.EncodeProjectPath(cwd)
	} else {
		// Encode the project path for database lookup
		*project = shared.EncodeProjectPath(*project)
	}

	// Open database
	database, err := db.Open(shared.DBPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	// Execute search
	matches, err := database.Search(query, *scope, *project, *limit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error searching: %v\n", err)
		os.Exit(1)
	}

	// Build result
	result := &db.SearchResult{
		Query:          query,
		Scope:          *scope,
		CurrentProject: *project,
		TotalMatches:   len(matches),
		Matches:        matches,
	}

	// Output results
	if *jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(result); err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
			os.Exit(1)
		}
	} else {
		printResults(result)
	}
}

func printResults(result *db.SearchResult) {
	fmt.Printf("Found %d conversation(s) matching \"%s\"\n\n", result.TotalMatches, result.Query)

	if len(result.Matches) == 0 {
		fmt.Println("No matches found.")
		return
	}

	for i, match := range result.Matches {
		fmt.Printf("%d. UUID: %s\n", i+1, match.UUID)
		fmt.Printf("   Project: %s\n", match.ProjectPath)

		// Parse and format created timestamp
		createdAt, err := time.Parse(time.RFC3339, match.CreatedAt)
		if err == nil {
			fmt.Printf("   Created: %s\n", createdAt.Format("Jan 2, 2006 at 3:04 PM"))
		} else {
			fmt.Printf("   Created: %s\n", match.CreatedAt)
		}

		fmt.Printf("   Messages: %d\n", match.MessageCount)
		fmt.Printf("   Summary: %s\n", match.Summary)
		fmt.Printf("   Relevance: %.2f\n\n", match.RelevanceScore)
	}
}
