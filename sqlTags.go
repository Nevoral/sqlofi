package sqlofi

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
	"time"
)

// TagOptions represents all possible SQLite column constraints and attributes
type TagOptions struct {
	Type string
	// Basic Constraints
	PrimaryKey    bool        // PRIMARY KEY constraint
	ForeignKey    *ForeignKey // FOREIGN KEY constraint with reference
	Unique        bool        // UNIQUE constraint
	NotNull       bool        // NOT NULL constraint
	AutoIncrement bool        // AUTOINCREMENT (only for INTEGER PRIMARY KEY)

	// Value Constraints
	Default string // DEFAULT value
	Check   string // CHECK constraint

	// Text-specific Attributes
	Collate string // COLLATION type (BINARY, NOCASE, RTRIM)

	// Generated Column
	Generated *GeneratedOption // GENERATED ALWAYS AS expression

	// Conflict Clause
	OnConflict string // ON CONFLICT clause (ROLLBACK, ABORT, FAIL, IGNORE, REPLACE)

	// Deferrable Constraints
	Deferrable        bool // DEFERRABLE
	InitiallyDeferred bool // INITIALLY DEFERRED

	// Index
	Index *IndexOptions // Creates an index on this column
}

func (t *TagOptions) Validate() error {
	// AUTOINCREMENT limitations
	if t.AutoIncrement {
		// Must be INTEGER PRIMARY KEY
		if !t.PrimaryKey || t.Type != "INTEGER" {
			return errors.New("AUTOINCREMENT is only allowed on INTEGER PRIMARY KEY columns")
		}
	}

	// Generated column limitations
	if t.Generated != nil {
		// Cannot have DEFAULT
		if t.Default != "" {
			return errors.New("generated column cannot have DEFAULT value")
		}
		// Cannot be PRIMARY KEY
		if t.PrimaryKey {
			return errors.New("generated column cannot be PRIMARY KEY")
		}
		// Cannot have UNIQUE constraint
		if t.Unique {
			return errors.New("generated column cannot be UNIQUE")
		}
		// Cannot be Foreign Key
		if t.ForeignKey != nil {
			return errors.New("generated column cannot be FOREIGN KEY")
		}
	}

	// PRIMARY KEY limitations
	if t.PrimaryKey {
		// Must be NOT NULL
		if !t.NotNull {
			t.NotNull = true // SQLite automatically makes PRIMARY KEY NOT NULL
		}
		// Cannot have DEFAULT
		if t.Default != "" {
			return errors.New("PRIMARY KEY column cannot have DEFAULT value")
		}
	}

	// FOREIGN KEY limitations
	if t.ForeignKey != nil {
		// Should match referenced column type
		// (This would require additional context about the referenced column)

		// Cannot be generated
		if t.Generated != nil {
			return errors.New("FOREIGN KEY column cannot be generated")
		}
	}

	// Index limitations
	if t.Index != nil {
		// Cannot create index on generated VIRTUAL column
		if t.Generated != nil && !t.Generated.Stored {
			return errors.New("cannot create index on VIRTUAL generated column")
		}
	}

	return nil
}

func (t *TagOptions) String() string {
	if t == nil {
		return ""
	}

	var parts []string

	// Add type if specified
	if t.Type != "" {
		parts = append(parts, t.Type)
	}

	// Add PRIMARY KEY and AUTOINCREMENT
	if t.PrimaryKey {
		if t.AutoIncrement {
			parts = append(parts, "PRIMARY KEY AUTOINCREMENT")
		} else {
			parts = append(parts, "PRIMARY KEY")
		}
	}

	// Add NOT NULL
	if t.NotNull {
		parts = append(parts, "NOT NULL")
	}

	// Add UNIQUE
	if t.Unique {
		parts = append(parts, "UNIQUE")
	}

	// Add DEFAULT
	if t.Default != "" {
		parts = append(parts, "DEFAULT "+t.Default)
	}

	// Add CHECK
	if t.Check != "" {
		parts = append(parts, "CHECK "+t.Check)
	}

	// Add COLLATE
	if t.Collate != "" {
		parts = append(parts, "COLLATE "+t.Collate)
	}

	// Add GENERATED
	if t.Generated != nil {
		parts = append(parts, t.Generated.String())
	}

	// Add ON CONFLICT
	if t.OnConflict != "" {
		parts = append(parts, "ON CONFLICT "+t.OnConflict)
	}

	// Add DEFERRABLE and INITIALLY DEFERRED
	if t.Deferrable {
		parts = append(parts, "DEFERRABLE")
		if t.InitiallyDeferred {
			parts = append(parts, "INITIALLY DEFERRED")
		}
	}

	// Note: ForeignKey and Index are typically handled separately
	// as they're usually added after the column definition

	return strings.Join(parts, " ")
}

// getSQLiteType returns the appropriate SQLite type for a Go type
func getSQLiteType(t reflect.Type) string {
	// Handle pointer types
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Handle sql.Null* types
	switch t {
	case reflect.TypeOf(sql.NullBool{}):
		return "BOOLEAN"
	case reflect.TypeOf(sql.NullInt16{}):
		return "INTEGER"
	case reflect.TypeOf(sql.NullInt32{}):
		return "INTEGER"
	case reflect.TypeOf(sql.NullInt64{}):
		return "INTEGER"
	case reflect.TypeOf(sql.NullFloat64{}):
		return "REAL"
	case reflect.TypeOf(sql.NullString{}):
		return "TEXT"
	case reflect.TypeOf(sql.NullTime{}):
		return "DATETIME"
	}

	// Handle basic Go types
	switch t.Kind() {
	case reflect.Bool:
		return "BOOLEAN"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "INTEGER"
	case reflect.Float32, reflect.Float64:
		return "REAL"
	case reflect.String:
		return "TEXT"
	}

	// Handle special types
	switch t {
	case reflect.TypeOf(time.Time{}):
		return "DATETIME"
	case reflect.TypeOf([]byte{}):
		return "BLOB"
	}

	// Default to TEXT for unknown types
	return "TEXT"
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
