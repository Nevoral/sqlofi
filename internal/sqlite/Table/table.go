package table

import (
	"fmt"
	"strings"

	reflectutil "github.com/Nevoral/sqlofi/internal/reflectUtil"
	check "github.com/Nevoral/sqlofi/internal/sqlite/Check"
	column "github.com/Nevoral/sqlofi/internal/sqlite/ColumnDef"
	expr "github.com/Nevoral/sqlofi/internal/sqlite/Expression"
	foreignkey "github.com/Nevoral/sqlofi/internal/sqlite/ForeignKey"
	primarykey "github.com/Nevoral/sqlofi/internal/sqlite/PrimaryKey"
	selectstmst "github.com/Nevoral/sqlofi/internal/sqlite/Select"
	unique "github.com/Nevoral/sqlofi/internal/sqlite/Unique"
	"github.com/Nevoral/sqlofi/internal/utils"
)

func NewTable(model any, foreignTablePtrs []any) *Table {
	return &Table{
		name:          reflectutil.GetStructName(model),
		model:         model,
		foreignTables: foreignTablePtrs,
	}
}

type Table struct {
	model         any
	foreignTables []any

	temporary   bool
	schema      string
	name        string
	ifNotExists bool
	selectSTMT  *selectstmst.Select

	constraints []string

	withoutRowID bool
	strict       bool
}

func (t *Table) Temporary() *Table {
	t.temporary = true
	return t
}

func (t *Table) IfNotExists() *Table {
	t.ifNotExists = true
	return t
}

func (t *Table) Schema(schemaName string) *Table {
	t.schema = schemaName
	return t
}

func (t *Table) Select(statement *selectstmst.Select) *Table {
	t.selectSTMT = statement
	return t
}

func (t *Table) WithouRowID() *Table {
	t.withoutRowID = true
	return t
}

func (t *Table) Strict() *Table {
	t.strict = true
	return t
}

func (t *Table) PrimaryKey(constraintName string, key *primarykey.TablePrimaryKey) *Table {
	if constraintName == "" {
		t.constraints = append(t.constraints, key.Build())
	} else {
		t.constraints = append(t.constraints,
			fmt.Sprintf("CONSTRAINT %s %s",
				constraintName,
				key.Build(),
			),
		)
	}
	return t
}

func (t *Table) Unique(constraintName string, unique *unique.TableUnique) *Table {
	if constraintName == "" {
		t.constraints = append(t.constraints, unique.Build())
	} else {
		t.constraints = append(t.constraints,
			fmt.Sprintf("CONSTRAINT %s %s",
				constraintName,
				unique.Build(),
			),
		)
	}
	return t
}

func (t *Table) Check(constraintName string, expression *expr.Expression) *Table {
	if constraintName == "" {
		t.constraints = append(t.constraints, check.NewCheck(expression).Build())
	} else {
		t.constraints = append(t.constraints,
			fmt.Sprintf("CONSTRAINT %s %s",
				constraintName,
				check.NewCheck(expression).Build(),
			),
		)
	}
	return t
}

func (t *Table) ForeignKey(constraintName string, key *foreignkey.References) *Table {
	if constraintName == "" {
		t.constraints = append(t.constraints, key.Build())
	} else {
		t.constraints = append(t.constraints,
			fmt.Sprintf("CONSTRAINT %s %s",
				constraintName,
				key.Build(),
			),
		)
	}
	return t
}

func (t *Table) Build() string {
	var (
		typeTable  = " TABLE"
		ifNotExist string
		schema     string
		body       string
	)
	if t.temporary {
		typeTable = fmt.Sprintf(" TEMP%s", typeTable)
	}

	if t.ifNotExists {
		ifNotExist = " IF NOT EXISTS"
	}

	if t.schema != "" {
		schema = t.schema + "."
	}

	if t.selectSTMT == nil {
		var options string
		if t.withoutRowID {
			options = " WITHOUT ROWID"
		}
		if t.strict {
			if options == "" {
				options = " STRICT"
			} else {
				options += ", STRICT"
			}
		}

		body = fmt.Sprintf("(\n%s%s\n)%s", t.buildColumnsDefinition(len(t.constraints) > 0), strings.Join(t.constraints, ",\n\t"), options)
	} else {
		body = fmt.Sprintf("AS %s", t.selectSTMT.Build())
	}

	return fmt.Sprintf("CREATE%s%s %s%s %s", typeTable, ifNotExist, schema, utils.ToSnakeCase(t.name), body)
}

func (t *Table) buildColumnsDefinition(existConstraints bool) string {
	var (
		result  = "\t"
		columns = t.getColumns()
	)

	for index, col := range columns {
		result += col.Build()
		if index+1 == len(columns) && !existConstraints {
			break
		}
		result += ",\n\t"
	}
	return result
}

func (t *Table) getColumns() []*column.Column {
	var columns []*column.Column
	for _, col := range reflectutil.GetStructFields(t.model) {
		ref := column.ParseStructField(t.foreignTables, col)
		columns = append(columns, ref)
	}
	return columns
}
