package sqlite

import (
	expr "github.com/Nevoral/sqlofi/internal/sqlite/Expression"
	selectstmt "github.com/Nevoral/sqlofi/internal/sqlite/Select"
)

type SelectType string

const (
	DISTINCT SelectType = "DISTINCT"
	ALL      SelectType = "ALL"
	NOTHING  SelectType = ""
)

func (s SelectType) String() string {
	return string(s)
}

// SELECT creates a new SELECT statement from a raw SQL string
func SELECT(selectType SelectType, columns ...*ResultColumn) *Select {
	covColumns := make([]*selectstmt.ResultColumn, len(columns))

	for i, column := range columns {
		covColumns[i] = column.ResultColumn
	}

	return &Select{
		Select: selectstmt.NewSelect(selectType.String(), covColumns),
	}
}

type Select struct {
	*selectstmt.Select
}

// From sets the FROM clause of the SELECT statement
func (s *Select) FROM(from *From) *Select {
	s.Select.From(from.From)
	return s
}

// Where sets the WHERE clause of the SELECT statement
func (s *Select) WHERE(condition *Expression) *Select {
	s.Select.Where(condition.Expression)
	return s
}

// GroupBy sets the GROUP BY clause of the SELECT statement
func (s *Select) GROUP_BY(expressions ...*Expression) *Select {
	var exprs []*expr.Expression
	for _, expr := range expressions {
		exprs = append(exprs, expr.Expression)
	}
	s.Select.GroupBy(exprs)
	return s
}

// Having sets the HAVING clause of the SELECT statement
func (s *Select) HAVING(condition *Expression) *Select {
	s.Select.Having(condition.Expression)
	return s
}

// OrderBy sets the ORDER BY clause of the SELECT statement
func (s *Select) ORDER_BY(orderBy ...*OrderBy) *Select {
	var orders []*selectstmt.OrderBy
	for _, order := range orderBy {
		orders = append(orders, order.OrderBy)
	}
	s.Select.OrderBy(orders...)
	return s
}

// Limit sets the LIMIT clause of the SELECT statement
func (s *Select) LIMIT(limit int) *Select {
	s.Select.Limit(limit)
	return s
}

// Offset sets the OFFSET clause of the SELECT statement
func (s *Select) OFFSET(offset int) *Select {
	s.Select.Offset(offset)
	return s
}

// NewExpressionColumn creates a new result column from an expression
func NewExpressionColumn(expression *Expression) *ResultColumn {
	return &ResultColumn{
		ResultColumn: selectstmt.NewExpressionColumn(expression.Expression),
	}
}

// NewExpressionColumnWithAlias creates a new result column from an expression with an alias
func NewExpressionColumnWithAlias(expression *Expression, alias string) *ResultColumn {
	return &ResultColumn{
		ResultColumn: selectstmt.NewExpressionColumnWithAlias(expression.Expression, alias),
	}
}

// NewWildcardColumn creates a new wildcard (*) result column
func NewWildcardColumn() *ResultColumn {
	return &ResultColumn{
		ResultColumn: selectstmt.NewWildcardColumn(),
	}
}

// NewTableWildcardColumn creates a new table wildcard (table.*) result column
func NewTableWildcardColumn(tableName string) *ResultColumn {
	return &ResultColumn{
		ResultColumn: selectstmt.NewTableWildcardColumn(tableName),
	}
}

// ResultColumn represents a column in the SELECT statement
type ResultColumn struct {
	*selectstmt.ResultColumn
}

// From represents a table or subquery in the FROM clause
type From struct {
	*selectstmt.From
}

// NewTableFrom creates a new FROM clause with a table
func NewTableFrom(tableName string) *From {
	return &From{
		From: selectstmt.NewTableFrom(tableName),
	}
}

// NewSubqueryFrom creates a new FROM clause with a subquery
func NewSubqueryFrom(subquery *Select) *From {
	return &From{
		From: selectstmt.NewSubqueryFrom(subquery.Select),
	}
}

// Alias sets an alias for the FROM clause
func (f *From) Alias(alias string) *From {
	f.From.Alias(alias)
	return f
}

// Join adds a join to the FROM clause
func (f *From) Join(join *Join) *From {
	f.From.Join(join.Join)
	return f
}

// JoinType represents the type of a JOIN
type JoinType string

const (
	INNER_JOIN       JoinType = "INNER JOIN"
	LEFT_JOIN        JoinType = "LEFT JOIN"
	LEFT_OUTER_JOIN  JoinType = "LEFT OUTER JOIN"
	RIGHT_JOIN       JoinType = "RIGHT JOIN"
	RIGHT_OUTER_JOIN JoinType = "RIGHT OUTER JOIN"
	FULL_JOIN        JoinType = "FULL JOIN"
	FULL_OUTER_JOIN  JoinType = "FULL OUTER JOIN"
	CROSS_JOIN       JoinType = "CROSS JOIN"
	NATURAL_JOIN     JoinType = "NATURAL JOIN"
)

func (j JoinType) String() string {
	return string(j)
}

// Join represents a JOIN clause
type Join struct {
	*selectstmt.Join
}

// NewTableJoin creates a new JOIN with a table
func NewTableJoin(joinType JoinType, tableName string) *Join {
	return &Join{
		Join: selectstmt.NewTableJoin(joinType.String(), tableName),
	}
}

// NewSubqueryJoin creates a new JOIN with a subquery
func NewSubqueryJoin(joinType JoinType, subquery *Select) *Join {
	return &Join{
		Join: selectstmt.NewSubqueryJoin(joinType.String(), subquery.Select),
	}
}

// Alias sets an alias for the JOIN
func (j *Join) Alias(alias string) *Join {
	j.Join.Alias(alias)
	return j
}

// On sets the ON condition for the JOIN
func (j *Join) On(condition *Expression) *Join {
	j.Join.On(condition.Expression)
	return j
}

// Using sets the USING columns for the JOIN
func (j *Join) Using(columns ...string) *Join {
	j.Join.Using(columns...)
	return j
}

// OrderDirection represents the direction of an ORDER BY clause
type OrderDirection string

const (
	ASC  OrderDirection = "ASC"
	DESC OrderDirection = "DESC"
)

func (d OrderDirection) String() string {
	return string(d)
}

// OrderBy represents an ORDER BY clause
type OrderBy struct {
	*selectstmt.OrderBy
}

// NewOrderBy creates a new ORDER BY clause
func NewOrderBy(expression *Expression, direction OrderDirection) *OrderBy {
	return &OrderBy{
		OrderBy: selectstmt.NewOrderBy(expression.Expression, direction.String()),
	}
}
