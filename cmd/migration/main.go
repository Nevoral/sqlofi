package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Nevoral/sqlofi/sqlite"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Initial schema models
type User struct {
	Id       int64          `sqlofi:"PRIMARY KEY AUTOINCREMENT"`
	Username string         `sqlofi:"NOT NULL UNIQUE"`
	Email    string         `sqlofi:"NOT NULL UNIQUE"`
	Password string         `sqlofi:"NOT NULL"`
	Created  sql.NullString `sqlofi:"DEFAULT (datetime('now'))"`
}

// Updated model with new fields
type UserV2 struct {
	Id        int64          `sqlofi:"PRIMARY KEY AUTOINCREMENT"`
	Username  string         `sqlofi:"NOT NULL UNIQUE"`
	Email     string         `sqlofi:"NOT NULL UNIQUE"`
	Password  string         `sqlofi:"NOT NULL"`
	FirstName sql.NullString `sqlofi:""`                   // New field
	LastName  sql.NullString `sqlofi:""`                   // New field
	IsActive  int            `sqlofi:"NOT NULL DEFAULT 1"` // New field
	LastLogin sql.NullString `sqlofi:""`                   // New field
	Created   sql.NullString `sqlofi:"DEFAULT (datetime('now'))"`
	Updated   sql.NullString `sqlofi:""` // New field
}

// Migration tracking table
type Migration struct {
	Id        int64  `sqlofi:"PRIMARY KEY AUTOINCREMENT"`
	Version   string `sqlofi:"NOT NULL UNIQUE"`
	AppliedAt string `sqlofi:"NOT NULL DEFAULT (datetime('now'))"`
}

func main() {
	dbPath := "migration_example.db"

	// Remove existing database if it exists (for demo purposes)
	os.Remove(dbPath)

	// Create a new database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Initial setup - create migrations table
	createMigrationsTable(db)

	// Apply initial schema (v1)
	applyMigrationV1(db)

	// Insert some sample data
	insertSampleUsers(db)

	// Apply schema update (v2)
	applyMigrationV2(db)

	// Show migration history
	showMigrationHistory(db)

	// Verify final schema
	verifyFinalSchema(db)
}

func createMigrationsTable(db *sql.DB) {
	fmt.Println("Creating migrations tracking table...")

	migrationTable := sqlite.CREATE_TABLE(Migration{}).IfNotExists()
	schema := sqlite.NewSchema("").Table(migrationTable)

	_, err := db.Exec(schema.Build())
	if err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	fmt.Println("Migrations table created")
}

func recordMigration(db *sql.DB, version string) error {
	_, err := db.Exec("INSERT INTO migration (Version) VALUES (?)", version)
	return err
}

func applyMigrationV1(db *sql.DB) {
	fmt.Println("\nApplying migration v1 - Initial schema...")

	// Check if migration was already applied
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM migration WHERE Version = ?", "v1").Scan(&count)
	if err == nil && count > 0 {
		fmt.Println("Migration v1 already applied, skipping...")
		return
	}

	// Create initial user table
	userTable := sqlite.CREATE_TABLE(User{}).IfNotExists()
	schema := sqlite.NewSchema("").
		Pragma(
			sqlite.ForeignKeys().ValueType("ON"),
		).
		Table(userTable)

	// Create the tables
	_, err = db.Exec(schema.Build())
	if err != nil {
		log.Fatalf("Failed to apply migration v1: %v", err)
	}

	// Record the migration
	err = recordMigration(db, "v1")
	if err != nil {
		log.Fatalf("Failed to record migration v1: %v", err)
	}

	fmt.Println("Migration v1 applied successfully")
}

func insertSampleUsers(db *sql.DB) {
	fmt.Println("\nInserting sample users...")

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to start transaction: %v", err)
	}

	// Insert some users
	users := []struct {
		username string
		email    string
		password string
	}{
		{"john_doe", "john@example.com", "hashed_password1"},
		{"jane_smith", "jane@example.com", "hashed_password2"},
		{"bob_jones", "bob@example.com", "hashed_password3"},
	}

	stmt, err := tx.Prepare("INSERT INTO user (Username, Email, Password) VALUES (?, ?, ?)")
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	for _, user := range users {
		_, err := stmt.Exec(user.username, user.email, user.password)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to insert user %s: %v", user.username, err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Println("Sample users inserted successfully")
}
func applyMigrationV2(db *sql.DB) {
	fmt.Println("\nApplying migration v2 - Adding new user fields...")

	// Check if migration was already applied
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM migration WHERE Version = ?", "v2").Scan(&count)
	if err == nil && count > 0 {
		fmt.Println("Migration v2 already applied, skipping...")
		return
	}

	// Start a transaction for the migration
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to start transaction: %v", err)
	}

	// Migration steps:

	// 1. Create a new temporary table with the new schema
	tempUserTable := sqlite.CREATE_TABLE(UserV2{}).Schema("temp")
	_, err = tx.Exec(tempUserTable.Build())
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to create temporary user table: %v", err)
	}

	// 2. Copy data from old table to new table
	_, err = tx.Exec(`
		INSERT INTO temp.user (Id, Username, Email, Password, Created)
		SELECT Id, Username, Email, Password, Created FROM main.user
	`)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to copy user data: %v", err)
	}

	// 3. Drop the old table
	_, err = tx.Exec("DROP TABLE main.user")
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to drop old user table: %v", err)
	}

	// 4. Rename the new table
	_, err = tx.Exec("ALTER TABLE temp.user RENAME TO user")
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to rename temporary user table: %v", err)
	}

	// 5. Create indexes on the new table
	_, err = tx.Exec("CREATE UNIQUE INDEX idx_user_username ON user(Username)")
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to create username index: %v", err)
	}

	_, err = tx.Exec("CREATE UNIQUE INDEX idx_user_email ON user(Email)")
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to create email index: %v", err)
	}

	// 6. Update the existing data with default values for new fields
	_, err = tx.Exec(`
		UPDATE user
		SET IsActive = 1,
		    Updated = datetime('now')
	`)
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to update user data: %v", err)
	}

	// Record the migration
	_, err = tx.Exec("INSERT INTO migration (Version) VALUES (?)", "v2")
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to record migration v2: %v", err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Fatalf("Failed to commit migration v2: %v", err)
	}

	fmt.Println("Migration v2 applied successfully")
}

func showMigrationHistory(db *sql.DB) {
	fmt.Println("\nMigration History:")
	fmt.Println("------------------")

	rows, err := db.Query("SELECT Version, AppliedAt FROM migration ORDER BY Id")
	if err != nil {
		log.Fatalf("Failed to query migration history: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version, appliedAt string
		if err := rows.Scan(&version, &appliedAt); err != nil {
			log.Fatalf("Failed to scan migration record: %v", err)
		}
		fmt.Printf("Version: %s, Applied: %s\n", version, appliedAt)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Error iterating migration rows: %v", err)
	}
}

func verifyFinalSchema(db *sql.DB) {
	fmt.Println("\nVerifying final schema:")
	fmt.Println("----------------------")

	// Show table structure
	rows, err := db.Query("PRAGMA table_info(user)")
	if err != nil {
		log.Fatalf("Failed to query table info: %v", err)
	}
	defer rows.Close()

	fmt.Println("User table columns:")
	fmt.Printf("%-3s %-15s %-10s %-5s %-15s %s\n", "Cid", "Name", "Type", "NotNull", "Default", "PK")
	fmt.Println(strings.Repeat("-", 70))

	for rows.Next() {
		var cid int
		var name, typeStr string
		var notNull int
		var dfltValue sql.NullString
		var pk int

		if err := rows.Scan(&cid, &name, &typeStr, &notNull, &dfltValue, &pk); err != nil {
			log.Fatalf("Failed to scan column info: %v", err)
		}

		defaultVal := "NULL"
		if dfltValue.Valid {
			defaultVal = dfltValue.String
		}

		fmt.Printf("%-3d %-15s %-10s %-5d %-15s %d\n", cid, name, typeStr, notNull, defaultVal, pk)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Error iterating column rows: %v", err)
	}

	// Show indexes
	fmt.Println("\nUser table indexes:")
	rows, err = db.Query("PRAGMA index_list(user)")
	if err != nil {
		log.Fatalf("Failed to query index list: %v", err)
	}
	defer rows.Close()

	fmt.Printf("%-3s %-20s %-7s %s\n", "Seq", "Name", "Unique", "Origin")
	fmt.Println(strings.Repeat("-", 50))

	for rows.Next() {
		var seq int
		var name, origin string
		var unique int

		if err := rows.Scan(&seq, &name, &unique, &origin); err != nil {
			log.Fatalf("Failed to scan index info: %v", err)
		}

		fmt.Printf("%-3d %-20s %-7d %s\n", seq, name, unique, origin)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Error iterating index rows: %v", err)
	}

	// Show sample data after migration
	fmt.Println("\nUser data after migration:")
	rows, err = db.Query(`
		SELECT Id, Username, Email, FirstName, LastName, IsActive, LastLogin, Created, Updated
		FROM user
	`)
	if err != nil {
		log.Fatalf("Failed to query user data: %v", err)
	}
	defer rows.Close()

	fmt.Printf("%-3s %-15s %-25s %-10s %-10s %-8s %-20s %-20s %-20s\n",
		"Id", "Username", "Email", "FirstName", "LastName", "IsActive", "LastLogin", "Created", "Updated")
	fmt.Println(strings.Repeat("-", 140))

	for rows.Next() {
		var id int64
		var username, email string
		var firstName, lastName, lastLogin, created, updated sql.NullString
		var isActive int

		if err := rows.Scan(&id, &username, &email, &firstName, &lastName, &isActive,
			&lastLogin, &created, &updated); err != nil {
			log.Fatalf("Failed to scan user data: %v", err)
		}

		// Format nullable fields
		firstNameStr := nullableString(firstName)
		lastNameStr := nullableString(lastName)
		lastLoginStr := nullableString(lastLogin)
		createdStr := nullableString(created)
		updatedStr := nullableString(updated)

		fmt.Printf("%-3d %-15s %-25s %-10s %-10s %-8d %-20s %-20s %-20s\n",
			id, username, email, firstNameStr, lastNameStr, isActive,
			lastLoginStr, createdStr, updatedStr)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("Error iterating user rows: %v", err)
	}

	fmt.Println("\nMigration test completed successfully!")
}

// Helper function to handle NULL strings
func nullableString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return "NULL"
}
