package selectstmt

import (
	"fmt"
	"strings"

	expr "github.com/Nevoral/sqlofi/internal/sqlite/Expression"
	"github.com/Nevoral/sqlofi/internal/utils"
)

// ResultColumnType represents the type of a result column
type ResultColumnType int

const (
	EXPRESSION     ResultColumnType = iota // expr [AS alias]
	WILDCARD                               // *
	TABLE_WILDCARD                         // table.*
)

// ResultColumn represents a column in the SELECT statement
type ResultColumn struct {
	columnType ResultColumnType
	expression *expr.Expression // Used for EXPRESSION type
	alias      string           // Optional alias for EXPRESSION type
	tableName  string           // Used for TABLE_WILDCARD type
}

// NewExpressionColumn creates a new result column from an expression
func NewExpressionColumn(expression *expr.Expression) *ResultColumn {
	return &ResultColumn{
		columnType: EXPRESSION,
		expression: expression,
	}
}

// NewExpressionColumnWithAlias creates a new result column from an expression with an alias
func NewExpressionColumnWithAlias(expression *expr.Expression, alias string) *ResultColumn {
	return &ResultColumn{
		columnType: EXPRESSION,
		expression: expression,
		alias:      alias,
	}
}

// NewWildcardColumn creates a new wildcard (*) result column
func NewWildcardColumn() *ResultColumn {
	return &ResultColumn{
		columnType: WILDCARD,
	}
}

// NewTableWildcardColumn creates a new table wildcard (table.*) result column
func NewTableWildcardColumn(tableName string) *ResultColumn {
	return &ResultColumn{
		columnType: TABLE_WILDCARD,
		tableName:  tableName,
	}
}

// Alias sets an alias for the result column
func (r *ResultColumn) Alias(alias string) *ResultColumn {
	if r.columnType != EXPRESSION {
		// Cannot set alias for wildcards
		return r
	}
	r.alias = alias
	return r
}

// Build returns the SQL representation of the result column
func (r *ResultColumn) Build() string {
	switch r.columnType {
	case EXPRESSION:
		if r.alias != "" {
			return fmt.Sprintf("%s AS %s", r.expression.Build(), r.alias)
		}
		return r.expression.Build()
	case WILDCARD:
		return "*"
	case TABLE_WILDCARD:
		return fmt.Sprintf("%s.*", utils.ToSnakeCase(r.tableName))
	default:
		return ""
	}
}

// From represents a table or subquery in the FROM clause
type From struct {
	tableName string
	alias     string
	subquery  *Select
	joins     []*Join
}

// NewTableFrom creates a new FROM clause with a table
func NewTableFrom(tableName string) *From {
	return &From{
		tableName: tableName,
	}
}

// NewSubqueryFrom creates a new FROM clause with a subquery
func NewSubqueryFrom(subquery *Select) *From {
	return &From{
		subquery: subquery,
	}
}

// Alias sets an alias for the FROM clause
func (f *From) Alias(alias string) *From {
	f.alias = alias
	return f
}

// Join adds a join to the FROM clause
func (f *From) Join(join *Join) *From {
	f.joins = append(f.joins, join)
	return f
}

// Build returns the SQL representation of the FROM clause
func (f *From) Build() string {
	var source string
	if f.tableName != "" {
		source = utils.ToSnakeCase(f.tableName)
	} else if f.subquery != nil {
		source = fmt.Sprintf("(%s)", f.subquery.Build())
	}

	if f.alias != "" {
		source = fmt.Sprintf("%s AS %s", source, f.alias)
	}

	if len(f.joins) == 0 {
		return source
	}

	var joins []string
	for _, join := range f.joins {
		joins = append(joins, join.Build())
	}

	return fmt.Sprintf("%s %s", source, strings.Join(joins, " "))
}

// Join represents a JOIN clause
type Join struct {
	joinType  string
	tableName string
	subquery  *Select
	alias     string
	on        *expr.Expression
	using     []string
}

// NewTableJoin creates a new JOIN with a table
func NewTableJoin(joinType string, tableName string) *Join {
	return &Join{
		joinType:  joinType,
		tableName: tableName,
	}
}

// NewSubqueryJoin creates a new JOIN with a subquery
func NewSubqueryJoin(joinType string, subquery *Select) *Join {
	return &Join{
		joinType: joinType,
		subquery: subquery,
	}
}

// Alias sets an alias for the JOIN
func (j *Join) Alias(alias string) *Join {
	j.alias = alias
	return j
}

// On sets the ON condition for the JOIN
func (j *Join) On(condition *expr.Expression) *Join {
	j.on = condition
	return j
}

// Using sets the USING columns for the JOIN
func (j *Join) Using(columns ...string) *Join {
	j.using = columns
	return j
}

// Build returns the SQL representation of the JOIN
func (j *Join) Build() string {
	var source string
	if j.tableName != "" {
		source = utils.ToSnakeCase(j.tableName)
	} else if j.subquery != nil {
		source = fmt.Sprintf("(%s)", j.subquery.Build())
	}

	if j.alias != "" {
		source = fmt.Sprintf("%s AS %s", source, j.alias)
	}

	var condition string
	if j.on != nil {
		condition = fmt.Sprintf(" ON %s", j.on.Build())
	} else if len(j.using) > 0 {
		// Convert column names to snake_case
		for i, col := range j.using {
			j.using[i] = utils.ToSnakeCase(col)
		}
		condition = fmt.Sprintf(" USING (%s)", strings.Join(j.using, ", "))
	}

	return fmt.Sprintf("%s %s%s", j.joinType, source, condition)
}

// Select represents a SELECT statement
type Select struct {
	selectType    string
	resultColumns []*ResultColumn
	from          *From
	where         *expr.Expression
	groupBy       []*expr.Expression
	having        *expr.Expression
	orderBy       []*OrderBy
	limit         int
	offset        int
	hasLimit      bool
	hasOffset     bool
	statement     string // Used for raw SQL statements
}

// OrderBy represents an ORDER BY clause
type OrderBy struct {
	expression *expr.Expression
	direction  string
}

// NewOrderBy creates a new ORDER BY clause
func NewOrderBy(expression *expr.Expression, direction string) *OrderBy {
	return &OrderBy{
		expression: expression,
		direction:  direction,
	}
}

// Build returns the SQL representation of the ORDER BY clause
func (o *OrderBy) Build() string {
	if o.direction == "" {
		return o.expression.Build()
	}
	return fmt.Sprintf("%s %s", o.expression.Build(), o.direction)
}

// NewSelect creates a new SELECT statement from a raw SQL string
func NewSelect(selectType string, columns []*ResultColumn) *Select {
	return &Select{
		resultColumns: columns,
	}
}

// From sets the FROM clause of the SELECT statement
func (s *Select) From(from *From) *Select {
	s.from = from
	return s
}

// Where sets the WHERE clause of the SELECT statement
func (s *Select) Where(condition *expr.Expression) *Select {
	s.where = condition
	return s
}

// GroupBy sets the GROUP BY clause of the SELECT statement
func (s *Select) GroupBy(expressions []*expr.Expression) *Select {
	s.groupBy = expressions
	return s
}

// Having sets the HAVING clause of the SELECT statement
func (s *Select) Having(condition *expr.Expression) *Select {
	s.having = condition
	return s
}

// OrderBy sets the ORDER BY clause of the SELECT statement
func (s *Select) OrderBy(orderBy ...*OrderBy) *Select {
	s.orderBy = orderBy
	return s
}

// Limit sets the LIMIT clause of the SELECT statement
func (s *Select) Limit(limit int) *Select {
	s.limit = limit
	s.hasLimit = true
	return s
}

// Offset sets the OFFSET clause of the SELECT statement
func (s *Select) Offset(offset int) *Select {
	s.offset = offset
	s.hasOffset = true
	return s
}

// Build returns the SQL representation of the SELECT statement
func (s *Select) Build() string {
	// If a raw statement was provided, return it
	if s.statement != "" {
		return s.statement
	}

	var parts []string

	// SELECT part
	var selectPart string
	if s.selectType != "" {
		selectPart = fmt.Sprintf("SELECT %s", s.selectType)
	} else {
		selectPart = "SELECT"
	}
	parts = append(parts, selectPart)

	// Result columns
	if len(s.resultColumns) == 0 {
		// Default to wildcard if no columns specified
		parts = append(parts, "*")
	} else {
		var columns []string
		for _, col := range s.resultColumns {
			columns = append(columns, col.Build())
		}
		parts = append(parts, strings.Join(columns, ", "))
	}

	// FROM clause
	if s.from != nil {
		parts = append(parts, fmt.Sprintf("FROM %s", s.from.Build()))
	}

	// WHERE clause
	if s.where != nil {
		parts = append(parts, fmt.Sprintf("WHERE %s", s.where.Build()))
	}

	// GROUP BY clause
	if len(s.groupBy) > 0 {
		var expressions []string
		for _, expr := range s.groupBy {
			expressions = append(expressions, expr.Build())
		}
		parts = append(parts, fmt.Sprintf("GROUP BY %s", strings.Join(expressions, ", ")))
	}

	// HAVING clause
	if s.having != nil {
		parts = append(parts, fmt.Sprintf("HAVING %s", s.having.Build()))
	}

	// ORDER BY clause
	if len(s.orderBy) > 0 {
		var orders []string
		for _, order := range s.orderBy {
			orders = append(orders, order.Build())
		}
		parts = append(parts, fmt.Sprintf("ORDER BY %s", strings.Join(orders, ", ")))
	}

	// LIMIT clause
	if s.hasLimit {
		parts = append(parts, fmt.Sprintf("LIMIT %d", s.limit))
	}

	// OFFSET clause
	if s.hasOffset {
		parts = append(parts, fmt.Sprintf("OFFSET %d", s.offset))
	}

	return strings.Join(parts, " ")
}

// Window is not implemented yet
func (s *Select) Window(expressions ...*expr.Expression) *Select {
	return s
}

// Union is not implemented yet
func (s *Select) Union(expressions ...*expr.Expression) *Select {
	return s
}

// UnionAll is not implemented yet
func (s *Select) UnionAll(expressions ...*expr.Expression) *Select {
	return s
}

// Intersect is not implemented yet
func (s *Select) Intersect(expressions ...*expr.Expression) *Select {
	return s
}

// Except is not implemented yet
func (s *Select) Except(expressions ...*expr.Expression) *Select {
	return s
}
