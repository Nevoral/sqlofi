# SQLOFI - Simple query syntax in go

SQLOFI is a simple Go library that makes database schema generation easy through struct tag annotations. This project was created as a learning exercise in Go, aiming to explore reflection, and SQL schema generation.

## Overview

SQLOFI allows you to define your database schema directly in your Go structs using tags. It parses these tags and generates appropriate SQL statements for creating tables, indexes, constraints, and relationships.

## Features

- Generate SQLite database schema from Go structs
- Support for common SQL constraints and features:
  - Primary keys
  - Foreign keys with different actions (CASCADE, SET NULL, etc.)
  - Indexes (including unique and partial indexes)
  - Generated/computed columns
  - Default values
  - Not null constraints
  - Auto-increment
- Automatic type mapping from Go types to SQLite types
- Simple API for setting up database schema

## Example Usage

```go
package main

import (
	"database/sql"
	"fmt"
	"github.com/yourusername/sqlofi"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID        int64  `sql:"primary key;autoincrement"`
	Username  string `sql:"NOT NULL;UNIQUE"`
	Email     string `sql:"NOT NULL;UNIQUE"`
	CreatedAt string `sql:"NOT NULL;DEFAULT:CURRENT_TIMESTAMP"`
}

type Post struct {
	ID      int64         `sql:"primary key;autoincrement"`
	Title   string        `sql:"NOT NULL"`
	Content string        `sql:"NOT NULL"`
	UserID  sql.NullInt64 `sql:"foreign key:user (id),ON DELETE CASCADE"`
	Likes   int           `sql:"DEFAULT:0"`
}

func main() {
	// Create schema from structs
	schema, err := sqlofi.CreateSchema(User{}, Post{})
	if err != nil {
		panic(err)
	}

	// Connect to database
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Set up database schema
	schema.SetUpDB(db)

	fmt.Println("Database schema created successfully")
}
```

## Tag Syntax

SQLOFI uses struct tags to define column properties:

```go
type MyStruct struct {
    ID       int64  `sql:"primary key;autoincrement"`
    Name     string `sql:"NOT NULL;UNIQUE"`
    ParentID int64  `sql:"foreign key:parent (id),ON DELETE CASCADE"`
    FullName string `sql:"GENERATED:Name || ' ' || LastName,STORED"`
}
```

Available tag options include:
- `primary key` - Makes the column a primary key
- `autoincrement` - Adds auto-increment (only for INTEGER PRIMARY KEY)
- `NOT NULL` - Adds NOT NULL constraint
- `UNIQUE` - Adds UNIQUE constraint
- `DEFAULT:value` - Sets default value
- `foreign key:table (column),action` - Creates foreign key reference
- `GENERATED:expression,STORED/VIRTUAL` - Creates computed column

## Project Status

This is a learning project and not intended for production use. It's a simple implementation to explore Go's capabilities for working with struct tags and database schemas.

## License

MIT

## Contributing

This is a personal learning project, but feel free to use it as inspiration for your own exploration of Go programming.
