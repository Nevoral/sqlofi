package foreignkey

import (
	"fmt"
	"slices"
	"strings"

	reflectutil "github.com/Nevoral/sqlofi/internal/reflectUtil"
	"github.com/Nevoral/sqlofi/internal/utils"
)

// ParseReference parses a reference from a string
// Expected format: "REFERENCES tableName (col, ...), ON DELETE ... ON UPDATE ... MATCH ... DEFERRABLE..."
func ParseReference(colName string, foreignTables []any, tag string) (*References, error) {
	tag = strings.TrimSpace(tag)

	// Check if the content starts with REFERENCES
	if !strings.HasPrefix(strings.ToUpper(tag), "REFERENCES") {
		tok := strings.Split(tag, " ")
		if tok != nil {
			return nil, fmt.Errorf("Error tag doesn't start with REFERENCES instead starts with %s", tok[0])
		}
		return nil, fmt.Errorf("Error tag doesn't start with REFERENCES and doesn't contains any space")
	}

	// Remove the REFERENCES prefix
	tag = strings.TrimSpace(tag[len("REFERENCES"):])

	refs := &References{
		columnsName: []string{colName},
	}

	var (
		index            int
		tablePart        string
		length           = len(tag)
		foreignTableName string
	)

	for {
		if tag[index] == ')' {
			tablePart = tag[:index+1]
			break
		}
		if index+3 < length {
			if tag[index:index+3] == " ON" {
				tablePart = tag[:index]
				break
			}
		}
		if index+6 < length {
			if tag[index:index+6] == " MATCH" {
				tablePart = tag[:index]
				break
			}
		}
		if index+8 < length {
			if tag[index:index+8] == " NOT DEF" {
				tablePart = tag[:index]
				break
			}
		}
		if index+8 < length {
			if tag[index:index+8] == " DEFERRA" {
				tablePart = tag[:index]
				break
			}
		}
		index++
	}
	tag = strings.TrimSpace(tag[index:])
	index = 0

	// Parse table name and columns
	if idx := strings.Index(tablePart, "("); idx != -1 {
		foreignTableName = strings.TrimSpace(tablePart[:idx])

		// Find the closing parenthesis
		endIdx := strings.LastIndex(tablePart, ")")
		if endIdx != -1 && endIdx > idx {
			// Extract columns
			colStr := tablePart[idx+1 : endIdx]
			cols := strings.Split(colStr, ",")
			for _, col := range cols {
				col = strings.TrimSpace(col)
				if col != "" {
					refs.foreignColumnsName = append(refs.foreignColumnsName, col)
				}
			}
		}
	} else {
		foreignTableName = strings.TrimSpace(tablePart)
	}

	for _, t := range foreignTables {
		if foreignTableName == reflectutil.GetStructName(t) {
			refs.foreignTable = t
		}
	}

	if refs.foreignTable == nil {
		return nil, fmt.Errorf("Error from provided structs as foreign tables none of them are matching %s", foreignTableName)
	}

	if tag != "" {
		// Match patterns: ON DELETE, ON UPDATE, MATCH, DEFERRABLE, NOT DEFERRABLE
		clauses := []string{}
		patterns := []string{"ON DELETE", "ON UPDATE", "MATCH", "DEFERRABLE", "NOT DEFERRABLE"}

		for index < len(tag) {
			found := false

			for _, pattern := range patterns {
				if index+len(pattern) <= len(tag) &&
					tag[index:index+len(pattern)] == pattern {
					// Found the index of a clause
					nextStart := index + len(pattern)

					// Find the index of the next clause
					nextClauseStart := len(tag)
					for _, nextPattern := range patterns {
						if idx := strings.Index(tag[nextStart:], nextPattern); idx != -1 {
							if nextStart+idx < nextClauseStart {
								nextClauseStart = nextStart + idx
							}
						}
					}

					clauses = append(clauses, strings.TrimSpace(tag[index:nextClauseStart]))
					index = nextClauseStart
					found = true
					break
				}
			}

			if !found {
				index++ // Move to next character if no pattern match at current position
			}

		}

		// Process each clause
		for _, clause := range clauses {
			switch {
			case strings.HasPrefix(clause, "ON DELETE"):
				actionStr := strings.TrimSpace(strings.TrimPrefix(clause, "ON DELETE"))
				if action, err := rowAction(actionStr); err == nil {
					refs.OnDelete(action)
				}

			case strings.HasPrefix(clause, "ON UPDATE"):
				actionStr := strings.TrimSpace(strings.TrimPrefix(clause, "ON UPDATE"))
				if action, err := rowAction(actionStr); err == nil {
					refs.OnUpdate(action)
				}

			case strings.HasPrefix(clause, "MATCH"):
				matchType := strings.TrimSpace(strings.TrimPrefix(clause, "MATCH"))
				refs.Match(matchType)

			case strings.HasPrefix(clause, "DEFERRABLE"):
				actionStr := strings.TrimSpace(strings.TrimPrefix(clause, "DEFERRABLE"))
				if action, err := defferableAction(actionStr); err == nil {
					refs.Deferrable(action)
				}

			case strings.HasPrefix(clause, "NOT DEFERRABLE"):
				actionStr := strings.TrimSpace(strings.TrimPrefix(clause, "NOT DEFERRABLE"))
				if action, err := defferableAction(actionStr); err == nil {
					refs.NotDeferrable(action)
				}
			}
		}
	}
	return refs, nil
}

func NewTableForeignTable(foreignTablePtr any, columns []string) *References {
	return &References{
		tableTypeReference: true,
		columnsName:        columns,
		foreignTable:       foreignTablePtr,
	}
}

func NewReferences(columName string, foreignTablePtr any, foreignColumn string) *References {
	return &References{
		tableTypeReference: false,
		columnsName:        []string{columName},
		foreignTable:       foreignTablePtr,
		foreignColumnsName: []string{foreignColumn},
	}
}

type References struct {
	tableTypeReference bool
	columnsName        []string

	foreignTable       any
	foreignColumnsName []string

	onDeleteVal      string
	onUpdateVal      string
	matchVal         string
	deferrableVal    *string
	notDeferrableVal *string
}

func (r *References) GetColumns() []string {
	return r.columnsName
}

func (r *References) OnDelete(action string) *References {
	r.onDeleteVal = action
	return r
}

func (r *References) OnUpdate(action string) *References {
	r.onUpdateVal = action
	return r
}

func (r *References) Match(name string) *References {
	r.matchVal = name
	return r
}

func (r *References) Deferrable(action string) *References {
	if r.notDeferrableVal != nil {
		return r
	}

	r.deferrableVal = &action
	return r
}

func (r *References) NotDeferrable(action string) *References {
	if r.deferrableVal != nil {
		return r
	}

	r.notDeferrableVal = &action
	return r
}

func (r *References) ForeighColumns(columns []string) *References {
	if !r.tableTypeReference {
		panic("Error: foreign key reference has to have only one")
	}
	if len(columns) != len(r.columnsName) {
		panic(fmt.Errorf("Error len of foreign table columns isn't same as len of columns"))
	}
	r.foreignColumnsName = columns
	return r
}

func (r *References) ForeignTable(foreignTablePtr any) *References {
	r.foreignTable = foreignTablePtr
	return r
}

func (r *References) Build() string {
	var (
		prefix         string
		colName        string
		actions        string
		tableName      = reflectutil.GetStructName(r.foreignTable)
		foreignColName = reflectutil.GetStructFieldsNames(r.foreignTable)
	)

	for _, col := range r.foreignColumnsName {
		if !slices.Contains(foreignColName, col) {
			panic(fmt.Errorf("column '%s' not found in foreign table '%s'", col, tableName))
		}
	}

	if r.tableTypeReference {
		prefix = fmt.Sprintf("FOREIGN KEY (%s) ", utils.Join(r.columnsName, ", "))
	}

	if len(r.foreignColumnsName) > 0 {
		colName = " ("
		colName += utils.Join(r.foreignColumnsName, ", ")
		colName += ")"
	}

	if r.onDeleteVal != "" {
		actions += fmt.Sprintf(" ON DELETE %s", r.onDeleteVal)
	}
	if r.onUpdateVal != "" {
		actions += fmt.Sprintf(" ON UPDATE %s", r.onUpdateVal)
	}
	if r.matchVal != "" {
		actions += fmt.Sprintf(" MATCH %s", r.matchVal)
	}
	if r.deferrableVal != nil {
		actions += " DEFERRABLE"
		if *r.deferrableVal != "" {
			actions += " " + *r.deferrableVal
		}
	}
	if r.notDeferrableVal != nil {
		actions += " NOT DEFERRABLE"
		if *r.notDeferrableVal != "" {
			actions += " " + *r.notDeferrableVal
		}
	}

	return fmt.Sprintf("%sREFERENCES %s%s%s", prefix, utils.ToSnakeCase(tableName), colName, actions)
}

func rowAction(value string) (string, error) {
	value = strings.TrimSpace(value)
	switch value {
	case "CASCADE", "SET NULL", "SET DEFAULT", "RESTRICT", "NO ACTION":
		return value, nil
	default:
		return "", fmt.Errorf("invalid RowAction: %s", value)
	}
}

func defferableAction(value string) (string, error) {
	value = strings.TrimSpace(value)
	switch value {
	case "INITIALLY DEFERRED", "INITIALLY IMMEDIATE", "":
		return value, nil
	default:
		return "", fmt.Errorf("invalid DeferrableAction: %s", value)
	}
}
