package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Nevoral/sqlofi/sqlite"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Define a struct for JSON data
type Metadata struct {
	Tags       []string          `json:"tags"`
	Properties map[string]string `json:"properties"`
	Version    int               `json:"version"`
}

// Document model with virtual columns and JSON
type Document struct {
	Id           int64          `sqlofi:"PRIMARY KEY AUTOINCREMENT"`
	Title        string         `sqlofi:"NOT NULL"`
	Content      string         `sqlofi:"NOT NULL"`
	MetadataJSON string         `sqlofi:"NOT NULL DEFAULT '{}'"`                                                                 // Stored as JSON text
	TagsCount    int            `sqlofi:"GENERATED ALWAYS AS (json_array_length(json_extract(MetadataJSON, '$.tags'))) VIRTUAL"` // Virtual column
	FirstTag     sql.NullString `sqlofi:"GENERATED ALWAYS AS (json_extract(MetadataJSON, '$.tags[0]')) STORED"`                  // Stored virtual column
	Created      string         `sqlofi:"NOT NULL DEFAULT (datetime('now'))"`
	Modified     string         `sqlofi:"NOT NULL DEFAULT (datetime('now'))"`
}

func main() {
	// Open database connection
	db, err := sql.Open("sqlite3", "advanced_features.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Enable JSON functions
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatalf("Failed to enable foreign keys: %v", err)
	}

	// Create document table
	documentTable := sqlite.CREATE_TABLE(Document{}).IfNotExists()

	schema := sqlite.NewSchema("advanced_features.db").
		Pragma(
			sqlite.JournalModeWAL(""),
			sqlite.SynchronousFull(""),
		).
		Table(
			documentTable,
		).
		Index(
			sqlite.CREATE_INDEX(&Document{}, "idx_document_tags_count").
				IfNotExists(),
			sqlite.CREATE_INDEX(&Document{}, "idx_document_first_tag").
				IfNotExists().
				Where(sqlite.NewExpression("FirstTag IS NOT NULL")),
		)

	// Execute schema
	_, err = db.Exec(schema.Build())
	if err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	fmt.Println("Schema created successfully")

	// Insert sample documents
	insertSampleDocuments(db)

	// Demo advanced queries using JSON and virtual columns
	demonstrateAdvancedFeatures(db)
}

func insertSampleDocuments(db *sql.DB) {
	fmt.Println("\nInserting sample documents...")

	// Prepare sample documents with JSON metadata
	documents := []struct {
		title    string
		content  string
		metadata Metadata
	}{
		{
			"Getting Started with SQLite JSON Features",
			"SQLite provides powerful JSON functions that allow you to store and query JSON data efficiently.",
			Metadata{
				Tags:       []string{"sqlite", "json", "tutorial"},
				Properties: map[string]string{"difficulty": "beginner", "type": "tutorial"},
				Version:    1,
			},
		},
		{
			"Advanced Query Techniques",
			"This document covers advanced query techniques including subqueries, CTEs, and window functions.",
			Metadata{
				Tags:       []string{"sqlite", "advanced", "query", "optimization"},
				Properties: map[string]string{"difficulty": "advanced", "type": "reference"},
				Version:    2,
			},
		},
		{
			"Virtual Columns in SQLite",
			"Virtual columns can be computed on-the-fly or stored, providing flexible ways to access derived data.",
			Metadata{
				Tags:       []string{"sqlite", "schema", "virtual-columns"},
				Properties: map[string]string{"difficulty": "intermediate", "type": "explanation"},
				Version:    1,
			},
		},
		{
			"Data Modeling Best Practices",
			"Learn how to design efficient and maintainable database schemas for your applications.",
			Metadata{
				Tags:       []string{"database", "design", "best-practices"},
				Properties: map[string]string{"difficulty": "intermediate", "type": "guide"},
				Version:    3,
			},
		},
		{
			"Empty Tags Example",
			"This document has no tags to demonstrate handling of empty arrays.",
			Metadata{
				Tags:       []string{},
				Properties: map[string]string{"type": "example"},
				Version:    1,
			},
		},
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}

	stmt, err := tx.Prepare(`
		INSERT INTO document (Title, Content, MetadataJSON)
		VALUES (?, ?, ?)
	`)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	for _, doc := range documents {
		// Convert metadata to JSON
		metadataBytes, err := json.Marshal(doc.metadata)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to marshal metadata: %v", err)
		}

		// Insert document
		_, err = stmt.Exec(doc.title, doc.content, string(metadataBytes))
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to insert document: %v", err)
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Println("Sample documents inserted successfully")
}

func demonstrateAdvancedFeatures(db *sql.DB) {
	fmt.Println("\nDemonstrating advanced SQLite features:")

	// 1. Query documents with their calculated fields
	fmt.Println("\n1. Documents with their calculated fields:")
	rows, err := db.Query(`
		SELECT
			Id,
			Title,
			TagsCount,
			FirstTag,
			json_extract(MetadataJSON, '$.version') as Version
		FROM document
		ORDER BY Id
	`)
	if err != nil {
		log.Fatalf("Failed to query documents: %v", err)
	}
	defer rows.Close()

	fmt.Printf("%-3s %-40s %-10s %-15s %-7s\n", "ID", "Title", "Tags Count", "First Tag", "Version")
	fmt.Println(strings.Repeat("-", 80))

	for rows.Next() {
		var id int64
		var title string
		var tagsCount int
		var firstTag sql.NullString
		var version int

		if err := rows.Scan(&id, &title, &tagsCount, &firstTag, &version); err != nil {
			log.Fatalf("Failed to scan document: %v", err)
		}

		firstTagStr := "NULL"
		if firstTag.Valid {
			firstTagStr = firstTag.String
		}

		fmt.Printf("%-3d %-40s %-10d %-15s %-7d\n", id, title, tagsCount, firstTagStr, version)
	}

	// 2. Query documents by tags with JSON path
	fmt.Println("\n2. Documents containing 'sqlite' tag (using JSON path):")
	rows, err = db.Query(`
		SELECT
			Id,
			Title,
			json_extract(MetadataJSON, '$.tags') as Tags
		FROM document
		WHERE json_extract(MetadataJSON, '$.tags') LIKE '%sqlite%'
		ORDER BY Id
	`)
	if err != nil {
		log.Fatalf("Failed to query documents by tag: %v", err)
	}
	defer rows.Close()

	fmt.Printf("%-3s %-40s %-30s\n", "ID", "Title", "Tags")
	fmt.Println(strings.Repeat("-", 75))

	for rows.Next() {
		var id int64
		var title, tags string

		if err := rows.Scan(&id, &title, &tags); err != nil {
			log.Fatalf("Failed to scan document: %v", err)
		}

		fmt.Printf("%-3d %-40s %-30s\n", id, title, tags)
	}

	// 3. Query documents by property with JSON path
	fmt.Println("\n3. Documents with 'intermediate' difficulty (using JSON path):")
	rows, err = db.Query(`
		SELECT
			Id,
			Title,
			json_extract(MetadataJSON, '$.properties.difficulty') as Difficulty,
			json_extract(MetadataJSON, '$.properties.type') as Type
		FROM document
		WHERE json_extract(MetadataJSON, '$.properties.difficulty') = 'intermediate'
		ORDER BY Id
	`)
	if err != nil {
		log.Fatalf("Failed to query documents by property: %v", err)
	}
	defer rows.Close()

	fmt.Printf("%-3s %-40s %-15s %-15s\n", "ID", "Title", "Difficulty", "Type")
	fmt.Println(strings.Repeat("-", 75))

	for rows.Next() {
		var id int64
		var title, difficulty, docType string

		if err := rows.Scan(&id, &title, &difficulty, &docType); err != nil {
			log.Fatalf("Failed to scan document: %v", err)
		}

		fmt.Printf("%-3d %-40s %-15s %-15s\n", id, title, difficulty, docType)
	}

	// 4. Update document metadata using JSON functions
	fmt.Println("\n4. Updating document metadata using JSON functions:")

	// Get first document
	var id int64
	var title string
	var metadataJSON string

	err = db.QueryRow("SELECT Id, Title, MetadataJSON FROM document LIMIT 1").Scan(&id, &title, &metadataJSON)
	if err != nil {
		log.Fatalf("Failed to get document: %v", err)
	}

	fmt.Printf("Updating document %d: %s\n", id, title)

	// Parse existing metadata
	var metadata Metadata
	err = json.Unmarshal([]byte(metadataJSON), &metadata)
	if err != nil {
		log.Fatalf("Failed to parse metadata: %v", err)
	}

	// Update metadata
	metadata.Tags = append(metadata.Tags, "updated")
	metadata.Properties["updated_at"] = time.Now().Format(time.RFC3339)
	metadata.Version++

	// Convert back to JSON
	updatedMetadataBytes, err := json.Marshal(metadata)
	if err != nil {
		log.Fatalf("Failed to marshal updated metadata: %v", err)
	}

	// Update the document
	_, err = db.Exec(`
		UPDATE document
		SET
			MetadataJSON = ?,
			Modified = datetime('now')
		WHERE Id = ?
	`, string(updatedMetadataBytes), id)
	if err != nil {
		log.Fatalf("Failed to update document: %v", err)
	}

	fmt.Println("Document updated successfully")

	// 5. Verify the update and check virtual columns
	fmt.Println("\n5. Verifying update and virtual columns:")

	err = db.QueryRow(`
		SELECT
			Title,
			TagsCount,
			FirstTag,
			json_extract(MetadataJSON, '$.version') as Version,
			json_extract(MetadataJSON, '$.properties.updated_at') as UpdatedAt
		FROM document
		WHERE Id = ?
	`, id).Scan(&title, &metadata.TagsCount, &firstTag, &metadata.Version, &metadata.Properties["updated_at"])
	if err != nil {
		log.Fatalf("Failed to get updated document: %v", err)
	}

	firstTagStr := "NULL"
	if firstTag.Valid {
		firstTagStr = firstTag.String
	}

	fmt.Printf("Title: %s\n", title)
	fmt.Printf("Tags Count: %d\n", metadata.TagsCount)
	fmt.Printf("First Tag: %s\n", firstTagStr)
	fmt.Printf("Version: %d\n", metadata.Version)
	fmt.Printf("Updated At: %s\n", metadata.Properties["updated_at"])

	// 6. Use JSON table function to explode tags into rows
	fmt.Println("\n6. Using JSON_EACH to list all tags across documents:")

	rows, err = db.Query(`
		WITH TagsArray AS (
			SELECT
				d.Id as DocId,
				d.Title,
				json_extract(d.MetadataJSON, '$.tags') as Tags
			FROM document d
			WHERE json_array_length(json_extract(d.MetadataJSON, '$.tags')) > 0
		)
		SELECT
			t.DocId,
			t.Title,
			json_each.value as Tag
		FROM TagsArray t, json_each(t.Tags)
		ORDER BY Tag, DocId
	`)
	if err != nil {
		log.Fatalf("Failed to query tags: %v", err)
	}
	defer rows.Close()

	fmt.Printf("%-3s %-40s %-15s\n", "ID", "Document", "Tag")
	fmt.Println(strings.Repeat("-", 60))

	for rows.Next() {
		var docId int64
		var docTitle, tag string

		if err := rows.Scan(&docId, &docTitle, &tag); err != nil {
			log.Fatalf("Failed to scan tag: %v", err)
		}

		fmt.Printf("%-3d %-40s %-15s\n", docId, docTitle, tag)
	}

	fmt.Println("\nAdvanced features demonstration completed")
}
