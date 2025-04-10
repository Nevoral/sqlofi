package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	index "github.com/Nevoral/sqlofi/internal/sqlite/Index"
	pragmas "github.com/Nevoral/sqlofi/internal/sqlite/Pragmas"
	table "github.com/Nevoral/sqlofi/internal/sqlite/Table"
)

func NewSchema(name string) *Schema {
	return &Schema{
		name: name,
	}
}

type Schema struct {
	name    string
	db      *sql.DB
	ctx     context.Context
	pragmas []*pragmas.Pragma
	tables  []*table.Table
	indexes []*index.Index
}

func (s *Schema) Pragma(pragmas ...*Pragma) *Schema {
	for _, prag := range pragmas {
		s.pragmas = append(s.pragmas, prag.Pragma)
	}
	return s
}

func (s *Schema) Table(tables ...*Table) *Schema {
	for _, tab := range tables {
		s.tables = append(s.tables, tab.Table)
	}
	return s
}

func (s *Schema) Index(indexes ...*Index) *Schema {
	for _, idx := range indexes {
		s.indexes = append(s.indexes, idx.Index)
	}
	return s
}

func (s *Schema) OpenDBConnection(driverName, dataSourceName string) (err error) {
	s.db, err = sql.Open(driverName, dataSourceName)
	return err
}

func (s *Schema) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "db down: %v", err)
		os.Exit(1)
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

// Close closes the database connection.
func (s *Schema) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Schema) SetUpDatabase() *Schema {
	for _, pragma := range s.pragmas {
		s.db.Exec(pragma.Build())
	}
	for _, table := range s.tables {
		s.db.Exec(table.Build())
	}
	for _, index := range s.indexes {
		s.db.Exec(index.Build())
	}
	return s
}

func (s *Schema) Build() string {
	var schema string
	for _, pragma := range s.pragmas {
		schema += fmt.Sprintf("%s;\n", pragma.Build())
	}
	schema += "\n"
	for _, table := range s.tables {
		schema += fmt.Sprintf("%s;\n\n", table.Build())
	}
	for _, index := range s.indexes {
		schema += fmt.Sprintf("%s;\n", index.Build())
	}
	return schema
}
