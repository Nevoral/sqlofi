package sqlofi

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

type SqlTable struct {
	Name        string
	TableConfig *TableOptions
	Columns     []*SqlColumn
	ForeignKeys []*ForeignKey
	Indexes     []*IndexOptions
}

// String generates the complete SQL CREATE TABLE statement
func (t *SqlTable) String() string {
	var b strings.Builder

	// Start CREATE TABLE
	fmt.Fprintf(&b, "CREATE TABLE IF NOT EXISTS %s (\n", t.Name)

	// Column definitions
	for i, col := range t.Columns {
		fmt.Fprintf(&b, "\t%s", col.String())
		if i < len(t.Columns)-1 || len(t.ForeignKeys) > 0 {
			b.WriteString(",")
		}
		b.WriteString("\n")
	}

	// Foreign key constraints
	for i, fk := range t.ForeignKeys {
		fmt.Fprintf(&b, "\t%s", fk.String())
		if i < len(t.ForeignKeys)-1 {
			b.WriteString(",")
		}
		b.WriteString("\n")
	}

	b.WriteString(");\n")

	// Create indexes
	for _, idx := range t.Indexes {
		b.WriteString(idx.String(t.Name))
		b.WriteString("\n")
	}

	return b.String()
}

func parseTag(field reflect.StructField) *SqlColumn {
	var (
		opts    TagOptions
		tag     = field.Tag.Get("sql")
		colName = toSnakeCase(field.Name)
	)
	if slices.Contains([]string{"", "-"}, tag) {
		return nil
	}

	opts.Type = getSQLiteType(field.Type)
	for _, p := range strings.Split(tag, ";") {
		key, content, found := strings.Cut(strings.TrimSpace(p), ":")

		if found {
			switch strings.ToUpper(key) {
			case "FOREIGN KEY":
				opts.ForeignKey = parseForeignKey(content)
				opts.ForeignKey.Column = colName
			case "DEFAULT":
				opts.Default = content
			case "CHECK":
				opts.Check = content
			case "COLLATE":
				opts.Collate = content
			case "GENERATED":
				opts.Generated = parseGenerated(content)
			case "ON CONFLICT":
				opts.OnConflict = content
			case "INDEX":
				opts.Index = parseIndex(false, content)
				opts.Index.Column = colName
			case "INDEX UNIQUE":
				opts.Index = parseIndex(true, content)
				opts.Index.Column = colName
			}
		} else {
			switch strings.ToUpper(key) {
			case "PRIMARY KEY":
				opts.PrimaryKey = true
			case "AUTOINCREMENT":
				opts.AutoIncrement = true
			case "NOT NULL":
				opts.NotNull = true
			case "UNIQUE":
				opts.Unique = true
			case "DEFERRABLE":
				opts.Deferrable = true
			case "INITIALLY DEFERRED":
				opts.InitiallyDeferred = true
			}
		}
	}
	if err := opts.Validate(); err != nil {
		panic(err)
	}
	return &SqlColumn{
		Name: colName,
		Tags: &opts,
	}
}

type SqlColumn struct {
	Name string
	Tags *TagOptions
}

// String generates the SQL column definition
func (c *SqlColumn) String() string {
	var parts []string
	parts = append(parts, c.Name)

	if c.Tags != nil {
		tagStr := c.Tags.String()
		if tagStr != "" {
			parts = append(parts, tagStr)
		}
	}

	return strings.Join(parts, " ")
}
