package column

import (
	"fmt"
	"reflect"
	"strings"

	check "github.com/Nevoral/sqlofi/internal/sqlite/Check"
	collate "github.com/Nevoral/sqlofi/internal/sqlite/Collate"
	defaultConstr "github.com/Nevoral/sqlofi/internal/sqlite/Default"
	expr "github.com/Nevoral/sqlofi/internal/sqlite/Expression"
	foreignkey "github.com/Nevoral/sqlofi/internal/sqlite/ForeignKey"
	generated "github.com/Nevoral/sqlofi/internal/sqlite/Generated"
	notnull "github.com/Nevoral/sqlofi/internal/sqlite/NotNull"
	primarykey "github.com/Nevoral/sqlofi/internal/sqlite/PrimaryKey"
	sortorder "github.com/Nevoral/sqlofi/internal/sqlite/SortOrder"
	types "github.com/Nevoral/sqlofi/internal/sqlite/Types"
	unique "github.com/Nevoral/sqlofi/internal/sqlite/Unique"
	"github.com/Nevoral/sqlofi/internal/utils"
)

func isConstraintKeyword(token string) bool {
	token = strings.ToUpper(token)
	return strings.Contains(string(CONSTRAINT), token) ||
		strings.Contains(string(PRIMARY_KEY), token) ||
		strings.Contains(string(NOT_NULL), token) ||
		strings.Contains(string(UNIQUE), token) ||
		strings.Contains(string(CHECK), token) ||
		strings.Contains(string(DEFAULT), token) ||
		strings.Contains(string(COLLATE), token) ||
		strings.Contains(string(REFERENCES), token) ||
		strings.Contains(string(GENERATED), token)
}

type constraintToken string

const (
	CONSTRAINT  constraintToken = "CONSTRAINT"
	PRIMARY_KEY constraintToken = "PRIMARY KEY"
	NOT_NULL    constraintToken = "NOT NULL"
	UNIQUE      constraintToken = "UNIQUE"
	CHECK       constraintToken = "CHECK"
	DEFAULT     constraintToken = "DEFAULT"
	COLLATE     constraintToken = "COLLATE"
	REFERENCES  constraintToken = "REFERENCES"
	GENERATED   constraintToken = "GENERATED"
	AS          constraintToken = "AS"
)

// ParseStructField parses a Go struct field into a Column definition
func ParseStructField(models []any, field reflect.StructField) *Column {
	tag := field.Tag.Get("sqlofi")
	if tag == "" || tag == "-" {
		return nil
	}

	colName := utils.ToSnakeCase(field.Name)

	// Get the default SQLite type based on Go type
	defaultType, _ := types.GetSQLiteType(field.Type)

	// Create a base column with default type
	column := NewColumn(colName, defaultType)
	column.models = models

	// Parse tag options
	if tag != "" && tag != "-" {
		column.parseColumnTag(tag)
	}

	return column
}

func NewColumn(name string, colType types.SQLiteType) *Column {
	return &Column{
		name:    name,
		colType: colType,
	}
}

type Column struct {
	name       string
	colType    types.SQLiteType
	constraint []string
	models     []any

	// Track which constraints have been added to prevent duplicates
	// and enforce constraint compatibility
	hasPrimaryKey    bool
	hasNotNull       bool
	hasUnique        bool
	hasDefault       bool
	hasCheck         bool
	hasCollate       bool
	hasForeignKey    bool
	hasGenerated     bool
	hasAutoincrement bool
}

// parseColumnTag parses the "sqlofi" tag and applies constraints to the column
func (c *Column) parseColumnTag(tag string) {
	if tag == "" {
		return
	}

	tokens := strings.Fields(tag)
	if len(tokens) == 0 {
		return
	}

	// Process the rest of the tokens
	for i := 0; i < len(tokens); i++ {
		var (
			token          = strings.ToUpper(tokens[i])
			constraintName string
		)

		// Handle named constraints with "CONSTRAINT" keyword
		if token == string(CONSTRAINT) && i+2 < len(tokens) {
			constraintName = tokens[i+1]
			token = strings.ToUpper(tokens[i+2])
			i += 2 // Move past constraint name and to the constraint type
		}

		switch {
		case strings.Contains(string(PRIMARY_KEY), token):
			i++
			// Similar logic as above but without constraint name
			sortOrder := sortorder.ASC
			conflict := ""
			autoincrement := false

			for j := i + 1; j < len(tokens); j++ {
				opt := strings.ToUpper(tokens[j])
				if opt == "ASC" {
					sortOrder = sortorder.ASC
					i = j
				} else if opt == "DESC" {
					sortOrder = sortorder.DESC
					i = j
				} else if opt == "AUTOINCREMENT" {
					autoincrement = true
					i = j
				} else if strings.HasPrefix(opt, "ON") && j+1 < len(tokens) {
					conflictStr := strings.ToUpper(tokens[j+1])
					conflict = parseConflictClause(conflictStr)
					i = j + 1
				} else if isConstraintKeyword(opt) {
					i = j - 1
					break
				}
			}

			c.PrimaryKey(constraintName, sortOrder, conflict, autoincrement)

		case strings.Contains(string(NOT_NULL), token):
			i++
			conflict := ""
			if i+2 < len(tokens) && strings.ToUpper(tokens[i+1]) == "ON" &&
				strings.ToUpper(tokens[i+2]) == "CONFLICT" && i+3 < len(tokens) {
				conflictStr := strings.ToUpper(tokens[i+3])
				conflict = parseConflictClause(conflictStr)
				i += 3
			}

			c.NotNull(constraintName, conflict)

		case strings.Contains(string(UNIQUE), token):
			conflict := ""
			if i+2 < len(tokens) && strings.ToUpper(tokens[i+1]) == "ON" &&
				strings.ToUpper(tokens[i+2]) == "CONFLICT" && i+3 < len(tokens) {
				conflictStr := strings.ToUpper(tokens[i+3])
				conflict = parseConflictClause(conflictStr)
				i += 3
			}

			c.Unique(constraintName, conflict)

		case strings.Contains(string(CHECK), token):
			checkExpr := findBalancedParentheses(tokens, i+1)
			if checkExpr != "" {
				expr := expr.NewExpression(checkExpr)
				c.Check(constraintName, expr)
				i += strings.Count(checkExpr, " ") + 1
			}

		case strings.Contains(string(DEFAULT), token):
			defaultVal := ""
			if i+1 < len(tokens) {
				if tokens[i+1] == "(" {
					defaultVal = findBalancedParentheses(tokens, i+1)
					i += strings.Count(defaultVal, " ") + 1
				} else {
					defaultVal = tokens[i+1]
					i++
				}

				c.Default(constraintName, defaultVal)
			}

		case strings.Contains(string(COLLATE), token):
			if i+1 < len(tokens) {
				c.Collate(constraintName, tokens[i+1])
				i++
			}

		case strings.Contains(string(REFERENCES), token):
			refStr := extractReferencesClause(tokens[i:])
			if refStr != "" {
				c.ForeignKey(constraintName, refStr)
				i += strings.Count(refStr, " ")
			}

		case strings.Contains(string(GENERATED), token):
			if i+3 < len(tokens) && strings.ToUpper(tokens[i+1]) == "ALWAYS" &&
				strings.ToUpper(tokens[i+2]) == "AS" {

				exprStr := findBalancedParentheses(tokens, i+3)
				if exprStr != "" {
					generatedExpr := expr.NewExpression(exprStr)

					storageType := generated.VIRTUAL // Default
					nextPos := i + 3 + strings.Count(exprStr, " ") + 1

					if nextPos < len(tokens) {
						storageStr := strings.ToUpper(tokens[nextPos])
						if storageStr == "STORED" {
							storageType = generated.STORED
							i = nextPos
						} else if storageStr == "VIRTUAL" {
							storageType = generated.VIRTUAL
							i = nextPos
						}
					}

					c.Generated(constraintName, true, generatedExpr, storageType)
					i += strings.Count(exprStr, " ") + 3
				}
			}
		}

	}
}

// PrimaryKey adds a PRIMARY KEY constraint to the column
func (c *Column) PrimaryKey(constraintName string, sortOrder sortorder.SortOrder, conflict string, autoincrement bool) *Column {
	// PRIMARY KEY columns are automatically NOT NULL
	c.hasNotNull = true

	// Cannot have both AUTOINCREMENT and non-INTEGER PRIMARY KEY
	if autoincrement && c.colType != types.INTEGER {
		panic("AUTOINCREMENT is only allowed on INTEGER PRIMARY KEY columns")
	}

	// Cannot have both PRIMARY KEY and GENERATED
	if c.hasGenerated {
		panic("GENERATED column cannot be PRIMARY KEY")
	}

	// Cannot have both DEFAULT and PRIMARY KEY
	if c.hasDefault {
		panic("PRIMARY KEY column cannot have DEFAULT value")
	}

	c.hasPrimaryKey = true
	c.hasAutoincrement = autoincrement

	// Add constraint
	if constraintName == "" {
		c.constraint = append(c.constraint, primarykey.NewColumnPrimaryKey(sortOrder).
			Conflict(conflict).
			Autoincrement(autoincrement).Build(),
		)
	} else {
		c.constraint = append(c.constraint,
			fmt.Sprintf("CONSTRAINT %s %s",
				constraintName,
				primarykey.NewColumnPrimaryKey(sortOrder).
					Conflict(conflict).
					Autoincrement(autoincrement).Build(),
			),
		)
	}
	return c
}

// NotNull adds a NOT NULL constraint to the column
func (c *Column) NotNull(constraintName string, conflict string) *Column {
	c.hasNotNull = true

	if constraintName == "" {
		c.constraint = append(c.constraint, notnull.NewNotNull(conflict))
	} else {
		c.constraint = append(c.constraint,
			fmt.Sprintf("CONSTRAINT %s %s",
				constraintName,
				notnull.NewNotNull(conflict),
			),
		)
	}
	return c
}

// Unique adds a UNIQUE constraint to the column
func (c *Column) Unique(constraintName string, conflict string) *Column {
	// Cannot have both UNIQUE and GENERATED
	if c.hasGenerated {
		panic("GENERATED column cannot be UNIQUE")
	}

	c.hasUnique = true

	if constraintName == "" {
		c.constraint = append(c.constraint, unique.NewColumnUnique(conflict))
	} else {
		c.constraint = append(c.constraint,
			fmt.Sprintf("CONSTRAINT %s %s",
				constraintName,
				unique.NewColumnUnique(conflict),
			),
		)
	}
	return c
}

// Check adds a CHECK constraint to the column
func (c *Column) Check(constraintName string, expr *expr.Expression) *Column {
	c.hasCheck = true

	if constraintName == "" {
		c.constraint = append(c.constraint, check.NewCheck(expr).Build())
	} else {
		c.constraint = append(c.constraint,
			fmt.Sprintf("CONSTRAINT %s %s",
				constraintName,
				check.NewCheck(expr).Build(),
			),
		)
	}
	return c
}

// Default adds a DEFAULT constraint to the column
func (c *Column) Default(constraintName string, content string) *Column {
	// Cannot have both DEFAULT and PRIMARY KEY
	if c.hasPrimaryKey {
		panic("PRIMARY KEY column cannot have DEFAULT value")
	}

	// Cannot have both DEFAULT and GENERATED
	if c.hasGenerated {
		panic("GENERATED column cannot have DEFAULT value")
	}

	c.hasDefault = true

	if constraintName == "" {
		c.constraint = append(c.constraint, defaultConstr.ParseDefault(content))
	} else {
		c.constraint = append(c.constraint,
			fmt.Sprintf("CONSTRAINT %s %s",
				constraintName,
				defaultConstr.ParseDefault(content),
			),
		)
	}
	return c
}

// Collate adds a COLLATE constraint to the column
func (c *Column) Collate(constraintName string, name string) *Column {
	c.hasCollate = true

	if constraintName == "" {
		c.constraint = append(c.constraint, collate.NewCollate(name))
	} else {
		c.constraint = append(c.constraint,
			fmt.Sprintf("CONSTRAINT %s %s",
				constraintName,
				collate.NewCollate(name),
			),
		)
	}
	return c
}

// ForeignKey adds a REFERENCES foreign key constraint to the column
func (c *Column) ForeignKey(constraintName string, content string) *Column {
	// Cannot have both FOREIGN KEY and GENERATED
	if c.hasGenerated {
		panic("GENERATED column cannot be FOREIGN KEY")
	}

	c.hasForeignKey = true

	ref, err := foreignkey.ParseReference(c.name, c.models, content)
	if err != nil {
		panic(err)
	}

	if constraintName == "" {
		c.constraint = append(c.constraint, ref.Build())
	} else {
		c.constraint = append(c.constraint,
			fmt.Sprintf("CONSTRAINT %s %s",
				constraintName,
				ref.Build(),
			),
		)
	}
	return c
}

// Generated adds a GENERATED ALWAYS constraint for computed columns
func (c *Column) Generated(constraintName string, always bool, expr *expr.Expression, storageType generated.StorageType) *Column {
	// Cannot have both GENERATED and PRIMARY KEY
	if c.hasPrimaryKey {
		panic("GENERATED column cannot be PRIMARY KEY")
	}

	// Cannot have both GENERATED and UNIQUE
	if c.hasUnique {
		panic("GENERATED column cannot be UNIQUE")
	}

	// Cannot have both GENERATED and DEFAULT
	if c.hasDefault {
		panic("GENERATED column cannot have DEFAULT value")
	}

	// Cannot have both GENERATED and FOREIGN KEY
	if c.hasForeignKey {
		panic("GENERATED column cannot be FOREIGN KEY")
	}

	c.hasGenerated = true

	if constraintName == "" {
		c.constraint = append(c.constraint, generated.NewGenerated(always, expr, storageType))
	} else {
		c.constraint = append(c.constraint,
			fmt.Sprintf("CONSTRAINT %s %s",
				constraintName,
				generated.NewGenerated(always, expr, storageType),
			),
		)
	}
	return c
}

func (c *Column) Build() string {
	constraintStr := ""
	if len(c.constraint) > 0 {
		constraintStr = " " + strings.Join(c.constraint, " ")
	}

	return fmt.Sprintf("%s%s%s", c.name, " "+c.colType.String(), constraintStr)
}

// Helper functions
func parseConflictClause(str string) string {
	switch strings.ToUpper(str) {
	case "ROLLBACK", "ABORT", "FAIL", "IGNORE", "REPLACE":
		return str
	default:
		return ""
	}
}

// Helper to find a balanced parentheses expression
func findBalancedParentheses(tokens []string, startIndex int) string {
	if startIndex >= len(tokens) || tokens[startIndex] != "(" {
		return ""
	}

	level := 0
	var builder strings.Builder

	for i := startIndex; i < len(tokens); i++ {
		token := tokens[i]

		// Count opening parentheses
		level += strings.Count(token, "(")

		// Count closing parentheses
		level -= strings.Count(token, ")")

		builder.WriteString(token)
		builder.WriteString(" ")

		if level == 0 {
			// We've found the matching closing parenthesis
			return strings.TrimSpace(builder.String())
		}
	}

	// No balanced parentheses found
	return ""
}

// Helper to extract the complete REFERENCES clause
func extractReferencesClause(tokens []string) string {
	if len(tokens) < 2 || strings.ToUpper(tokens[0]) != "REFERENCES" {
		return ""
	}

	var builder strings.Builder
	builder.WriteString(tokens[0]) // REFERENCES

	inParentheses := false

	for i := 1; i < len(tokens); i++ {
		token := tokens[i]

		if strings.HasPrefix(token, "(") {
			inParentheses = true
		} else if strings.HasSuffix(token, ")") {
			inParentheses = false
		} else if !inParentheses && isConstraintKeyword(token) {
			// Stop at the beginning of a new constraint keyword
			// except for the ones that could be part of a foreign key clause
			break
		}

		builder.WriteString(" ")
		builder.WriteString(token)
	}

	return builder.String()
}
