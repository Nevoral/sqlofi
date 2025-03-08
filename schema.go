package sqlofi

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
)

// GenerateSchema generates SQL schema from structs
func CreateSchema(models ...interface{}) (*SQLSchema, error) {
	schema := &SQLSchema{}

	for _, str := range models {
		t := reflect.TypeOf(str)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct {
			return nil, fmt.Errorf("input must be a struct, got %v", t.Kind())
		}

		table := SqlTable{
			Name:        toSnakeCase(t.Name()),
			Columns:     []*SqlColumn{},
			ForeignKeys: []*ForeignKey{},
			Indexes:     []*IndexOptions{},
		}
		indexMap := make(map[string]*IndexOptions)

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			col := parseTag(field)
			if col == nil {
				continue
			}

			if col.Tags.ForeignKey != nil {
				table.ForeignKeys = append(table.ForeignKeys, col.Tags.ForeignKey)
			}

			if col.Tags.Index != nil {
				indexName := col.Tags.Index.Name
				if existing, exists := indexMap[indexName]; exists {
					existing.Column = fmt.Sprintf("%s, %s", existing.Column, col.Name)

					existing.Unique = existing.Unique || col.Tags.Index.Unique
					if col.Tags.Index.Where != "" || existing.Where != "" {
						existing.Where = col.Tags.Index.Where
					}
				} else {
					indexMap[indexName] = col.Tags.Index
				}
			}
			table.Columns = append(table.Columns, col)
		}
		for _, indexes := range indexMap {
			table.Indexes = append(table.Indexes, indexes)
		}
		schema.Tables = append(schema.Tables, &table)
	}

	return schema, nil
}

type SQLSchema struct {
	Tables []*SqlTable
}

func (s *SQLSchema) SetUpDB(db *sql.DB) {
	for _, table := range s.Tables {
		if _, err := db.Exec(table.String()); err != nil {
			fmt.Fprintf(os.Stderr, "Error %v", err)
			os.Exit(1)
		}

		if table.Indexes != nil {
			for _, index := range table.Indexes {
				if _, err := db.Exec(index.String(table.Name)); err != nil {
					fmt.Fprintf(os.Stderr, "Error %v", err)
					os.Exit(1)
				}
			}
		}

	}
}
