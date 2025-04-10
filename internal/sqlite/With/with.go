package with

import (
	"fmt"
	"strings"

	reflectutil "github.com/Nevoral/sqlofi/internal/reflectUtil"
	selectstmst "github.com/Nevoral/sqlofi/internal/sqlite/Select"
	"github.com/Nevoral/sqlofi/internal/utils"
)

func NewWithClause(recursive bool, tableExpressions ...*TableExpresions) string {
	withClause := "WITH "
	if recursive {
		withClause += "RECURSIVE "
	}

	tableExprs := make([]string, 0, len(tableExpressions))
	for _, tableExpr := range tableExpressions {
		tableExprs = append(tableExprs, tableExpr.Build())
	}

	return withClause + strings.Join(tableExprs, ", ")
}

func NewTableExpresion(table any, columns ...string) *TableExpresions {
	return &TableExpresions{
		table:        table,
		columns:      columns,
		materialized: 0, // Default is neither materialized nor not materialized
	}
}

type TableExpresions struct {
	table        any
	columns      []string
	materialized int8 // -1: NOT MATERIALIZED, 0: default, 1: MATERIALIZED
	selectStmt   *selectstmst.Select
}

// Materialized sets the table expression as MATERIALIZED
func (t *TableExpresions) Materialized() *TableExpresions {
	t.materialized = 1
	return t
}

// NotMaterialized sets the table expression as NOT MATERIALIZED
func (t *TableExpresions) NotMaterialized() *TableExpresions {
	t.materialized = -1
	return t
}

// WithSelect specifies the SELECT statement for this table expression
func (t *TableExpresions) WithSelect(stmt *selectstmst.Select) *TableExpresions {
	t.selectStmt = stmt
	return t
}

func (t *TableExpresions) Build() string {
	tableName := reflectutil.GetStructName(t.table)

	// Build the column list if specified
	columnList := ""
	if len(t.columns) > 0 {
		columnList = "(" + utils.Join(t.columns, ", ") + ")"
	}

	// Build the materialized clause
	materialized := ""
	switch t.materialized {
	case 1:
		materialized = " MATERIALIZED"
	case -1:
		materialized = " NOT MATERIALIZED"
	}

	// Build the select statement
	if t.selectStmt == nil {
		return fmt.Sprintf("%s%s AS%s", utils.ToSnakeCase(tableName), columnList, materialized)
	}

	// We have a SELECT statement
	selectSQL := t.selectStmt.Build()

	return fmt.Sprintf("%s%s AS%s (%s)", utils.ToSnakeCase(tableName), columnList, materialized, selectSQL)
}
