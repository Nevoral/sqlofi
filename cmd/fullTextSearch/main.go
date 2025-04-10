package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Nevoral/sqlofi/sqlite"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Article model for regular table
type Article struct {
	Id        int64          `sqlofi:"PRIMARY KEY AUTOINCREMENT"`
	Title     string         `sqlofi:"NOT NULL"`
	Author    string         `sqlofi:"NOT NULL"`
	Content   string         `sqlofi:"NOT NULL"`
	Published sql.NullString `sqlofi:"DEFAULT (datetime('now'))"`
	Tags      sql.NullString `sqlofi:""`
}

func main() {
	// Open database
	db, err := sql.Open("sqlite3", "fts_example.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create regular articles table
	articleTable := sqlite.CREATE_TABLE(Article{}).IfNotExists()

	schema := sqlite.NewSchema("fts_example.db").
		Pragma(
			sqlite.ForeignKeys().ValueType("ON"),
			sqlite.JournalModeWAL(""),
		).
		Table(
			articleTable,
		)

	// Build and execute schema
	schemaSQL := schema.Build()
	_, err = db.Exec(schemaSQL)
	if err != nil {
		log.Fatalf("Failed to execute schema: %v", err)
	}

	// Create FTS virtual table (this needs manual SQL since virtual tables have special syntax)
	_, err = db.Exec(`
		CREATE VIRTUAL TABLE IF NOT EXISTS article_fts USING fts5(
			title,
			author,
			content,
			tags,
			content='article',
			tokenize='porter unicode61'
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create FTS virtual table: %v", err)
	}

	// Create triggers to keep FTS index in sync with the article table
	_, err = db.Exec(`
		-- Insert trigger
		CREATE TRIGGER IF NOT EXISTS article_ai AFTER INSERT ON article BEGIN
			INSERT INTO article_fts(rowid, title, author, content, tags)
			VALUES (new.Id, new.Title, new.Author, new.Content, new.Tags);
		END;

		-- Delete trigger
		CREATE TRIGGER IF NOT EXISTS article_ad AFTER DELETE ON article BEGIN
			INSERT INTO article_fts(article_fts, rowid, title, author, content, tags)
			VALUES('delete', old.Id, old.Title, old.Author, old.Content, old.Tags);
		END;

		-- Update trigger
		CREATE TRIGGER IF NOT EXISTS article_au AFTER UPDATE ON article BEGIN
			INSERT INTO article_fts(article_fts, rowid, title, author, content, tags)
			VALUES('delete', old.Id, old.Title, old.Author, old.Content, old.Tags);
			INSERT INTO article_fts(rowid, title, author, content, tags)
			VALUES (new.Id, new.Title, new.Author, new.Content, new.Tags);
		END;
	`)
	if err != nil {
		log.Fatalf("Failed to create FTS triggers: %v", err)
	}

	fmt.Println("Schema created successfully")

	// Insert sample articles
	sampleArticles := []struct {
		title   string
		author  string
		content string
		tags    string
	}{
		{
			"Getting Started with SQLite in Go",
			"Jane Smith",
			"SQLite is a lightweight, embedded relational database that is perfect for small to medium-sized applications. In this article, we'll explore how to use SQLite with Go.",
			"sqlite,golang,database,tutorial",
		},
		{
			"Advanced SQLite Features",
			"John Doe",
			"SQLite offers many advanced features like JSON support, window functions, and full-text search. This article dives deep into these capabilities and how to leverage them in your applications.",
			"sqlite,advanced,fts,json",
		},
		{
			"Performance Tuning for Go Applications",
			"Jane Smith",
			"Optimizing your Go applications for performance is crucial. This article covers profiling, benchmarking, and various techniques to make your Go code run faster.",
			"golang,performance,optimization",
		},
		{
			"Building RESTful APIs with Go",
			"Bob Johnson",
			"This tutorial walks through creating a RESTful API using Go and SQLite as the database. We'll cover routing, middleware, authentication, and best practices.",
			"golang,api,rest,web",
		},
		{
			"SQLite vs PostgreSQL: When to Use Which",
			"Alice Williams",
			"Choosing the right database for your project is important. This comparison of SQLite and PostgreSQL will help you understand the strengths and weaknesses of each.",
			"sqlite,postgresql,comparison,database",
		},
	}

	// Use a transaction for inserting sample data
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO article (Title, Author, Content, Tags)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	for _, article := range sampleArticles {
		_, err := stmt.Exec(article.title, article.author, article.content, article.tags)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to insert article: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Println("Sample articles inserted successfully")

	// Demonstrate FTS searches
	demonstrateFTSQueries(db)
}

func demonstrateFTSQueries(db *sql.DB) {
	fmt.Println("\nFull-Text Search Examples:")

	// Define some search queries to demonstrate
	searches := []struct {
		description string
		query       string
		params      []interface{}
	}{
		{
			"1. Basic search for 'sqlite':",
			`SELECT
				a.Id,
				a.Title,
				a.Author,
				snippet(article_fts, 0, '<b>', '</b>', '...', 10) as ContentSnippet
			FROM article_fts
			JOIN article a ON article_fts.rowid = a.Id
			WHERE article_fts MATCH ?
			ORDER BY rank`,
			[]interface{}{"sqlite"},
		},
		{
			"2. Search for articles by 'Jane Smith':",
			`SELECT
				a.Id,
				a.Title,
				snippet(article_fts, 2, '<b>', '</b>', '...', 10) as ContentSnippet
			FROM article_fts
			JOIN article a ON article_fts.rowid = a.Id
			WHERE article_fts MATCH ?
			ORDER BY rank`,
			[]interface{}{"author:\"Jane Smith\""},
		},
		{
			"3. Combined search (articles about performance in Go):",
			`SELECT
				a.Id,
				a.Title,
				a.Author,
				snippet(article_fts, 2, '<b>', '</b>', '...', 10) as ContentSnippet
			FROM article_fts
			JOIN article a ON article_fts.rowid = a.Id
			WHERE article_fts MATCH ?
			ORDER BY rank`,
			[]interface{}{"performance AND golang"},
		},
		{
			"4. Phrase search:",
			`SELECT
				a.Id,
				a.Title,
				a.Author,
				snippet(article_fts, 2, '<b>', '</b>', '...', 10) as ContentSnippet
			FROM article_fts
			JOIN article a ON article_fts.rowid = a.Id
			WHERE article_fts MATCH ?
			ORDER BY rank`,
			[]interface{}{`"restful api"`},
		},
		{
			"5. Tag search:",
			`SELECT
				a.Id,
				a.Title,
				a.Author,
				a.Tags
			FROM article_fts
			JOIN article a ON article_fts.rowid = a.Id
			WHERE article_fts MATCH ?
			ORDER BY rank`,
			[]interface{}{"tags:tutorial OR tags:comparison"},
		},
	}

	// Execute each search query
	for _, search := range searches {
		fmt.Println("\n" + search.description)
		fmt.Println(strings.Repeat("-", len(search.description)))

		rows, err := db.Query(search.query, search.params...)
		if err != nil {
			fmt.Printf("Error executing query: %v\n", err)
			continue
		}

		// Get column names
		cols, err := rows.Columns()
		if err != nil {
			fmt.Printf("Error getting columns: %v\n", err)
			rows.Close()
			continue
		}

		// Create a slice of interface{} to hold the values
		values := make([]interface{}, len(cols))
		scanArgs := make([]interface{}, len(cols))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		// Print results
		resultsFound := false
		for rows.Next() {
			resultsFound = true
			err = rows.Scan(scanArgs...)
			if err != nil {
				fmt.Printf("Error scanning row: %v\n", err)
				continue
			}

			// Print each column
			for i, col := range cols {
				val := values[i]
				if val == nil {
					fmt.Printf("%s: [NULL]\n", col)
				} else {
					switch v := val.(type) {
					case []byte:
						fmt.Printf("%s: %s\n", col, string(v))
					default:
						fmt.Printf("%s: %v\n", col, v)
					}
				}
			}
			fmt.Println()
		}
		rows.Close()

		if !resultsFound {
			fmt.Println("No results found")
		}
	}

	// Demonstrate a more advanced use case: highlighting matches in context
	fmt.Println("\n6. Advanced example - Rank results by custom scoring and highlight matches:")
	query := `
		SELECT
			a.Id,
			a.Title,
			a.Author,
			highlight(article_fts, 0, '<mark>', '</mark>') as HighlightedTitle,
			highlight(article_fts, 2, '<mark>', '</mark>') as HighlightedContent,
			-- Custom ranking: title matches are more important than content matches
			(
				bm25(article_fts) +
				(CASE WHEN title MATCH ? THEN 10 ELSE 0 END) +
				(CASE WHEN tags MATCH ? THEN 5 ELSE 0 END)
			) as CustomScore
		FROM article_fts
		JOIN article a ON article_fts.rowid = a.Id
		WHERE article_fts MATCH ?
		ORDER BY CustomScore DESC
	`

	searchTerm := "sqlite"
	rows, err := db.Query(query, searchTerm, searchTerm, searchTerm)
	if err != nil {
		fmt.Printf("Error executing advanced query: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Printf("Search results for '%s' (ranked by custom score):\n", searchTerm)
	fmt.Println(strings.Repeat("-", 50))

	for rows.Next() {
		var id int64
		var title, author, highlightedTitle, highlightedContent string
		var customScore float64

		err := rows.Scan(&id, &title, &author, &highlightedTitle, &highlightedContent, &customScore)
		if err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			continue
		}

		fmt.Printf("ID: %d | Score: %.2f\n", id, customScore)
		fmt.Printf("Title: %s\n", highlightedTitle)
		fmt.Printf("Author: %s\n", author)
		fmt.Printf("Content Excerpt: %s\n",
			highlightedContent[:min(len(highlightedContent), 150)]+"...")
		fmt.Println(strings.Repeat("-", 50))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
