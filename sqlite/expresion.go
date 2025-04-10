package sqlite

import (
	"fmt"
	"strings"

	expr "github.com/Nevoral/sqlofi/internal/sqlite/Expression"
	selectstmt "github.com/Nevoral/sqlofi/internal/sqlite/Select"
	types "github.com/Nevoral/sqlofi/internal/sqlite/Types"
)

func Expr[T types.LiteralValue | types.DbPath | types.BindingParameter | string | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | []byte | float32 | float64 | bool](expression T) *Expression {
	switch ex := any(expression).(type) {
	case types.LiteralValue:
		return &Expression{Expression: expr.NewExpression(ex.String())}
	case types.DbPath:
		return &Expression{Expression: expr.NewExpression(ex.StringColumn())}
	case types.BindingParameter:
		return &Expression{Expression: expr.NewExpression(ex.String())}
	case string:
		return &Expression{Expression: expr.NewExpression("\"" + ex + "\"")}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return &Expression{Expression: expr.NewExpression(fmt.Sprintf("%d", ex))}
	case []byte:
		return &Expression{Expression: expr.NewExpression("\"" + string(ex) + "\"")}
	case float32, float64:
		return &Expression{Expression: expr.NewExpression(fmt.Sprintf("%f", ex))}
	case bool:
		if ex {
			return &Expression{Expression: expr.NewExpression("TRUE")}
		}
		return &Expression{Expression: expr.NewExpression("FALSE")}
	}
	return &Expression{Expression: expr.NewExpression("")}
}

type Expression struct {
	*expr.Expression
}

func Expressions(expressions ...*Expression) string {
	var expr string
	for idx, expression := range expressions {
		if expression == nil {
			panic("nil expression")
		}
		if idx > 0 {
			expr += ", "
		}
		expr += expression.Build()
	}
	return fmt.Sprintf("(%s)", expr)
}

func CAST(expression *Expression, typeName types.SQLiteType) string {
	return fmt.Sprintf("CAST(%s AS %s)", expression.Build(), typeName)
}

func COLLATE(expression *Expression, collation string) string {
	return fmt.Sprintf("%s COLLATE %s", expression.Build(), collation)
}

func NOT_LIKE(firstExpr *Expression, likeExpr *Expression, escapeExpr *Expression) string {
	if escapeExpr == nil {
		return fmt.Sprintf("%s NOT LIKE %s", firstExpr.Build(), likeExpr.Build())
	}
	return fmt.Sprintf("%s NOT LIKE %s ESCAPE %s", firstExpr.Build(), likeExpr.Build(), escapeExpr.Build())
}

func LIKE(firstExpr *Expression, likeExpr *Expression, escapeExpr *Expression) string {
	if escapeExpr == nil {
		return fmt.Sprintf("%s LIKE %s", firstExpr.Build(), likeExpr.Build())
	}
	return fmt.Sprintf("%s LIKE %s ESCAPE %s", firstExpr.Build(), likeExpr.Build(), escapeExpr.Build())
}

func NOT_GLOB(firstExpr *Expression, globExpr *Expression) string {
	return fmt.Sprintf("%s NOT GLOB %s", firstExpr.Build(), globExpr.Build())
}

func GLOB(firstExpr *Expression, globExpr *Expression) string {
	return fmt.Sprintf("%s GLOB %s", firstExpr.Build(), globExpr.Build())
}

func NOT_REGEXP(firstExpr *Expression, regexpExpr *Expression) string {
	return fmt.Sprintf("%s NOT REGEXP %s", firstExpr.Build(), regexpExpr.Build())
}

func REGEXP(firstExpr *Expression, regexpExpr *Expression) string {
	return fmt.Sprintf("%s REGEXP %s", firstExpr.Build(), regexpExpr.Build())
}

func NOT_MATCH(firstExpr *Expression, matchExpr *Expression) string {
	return fmt.Sprintf("%s NOT MATCH %s", firstExpr.Build(), matchExpr.Build())
}

func MATCH(firstExpr *Expression, matchExpr *Expression) string {
	return fmt.Sprintf("%s MATCH %s", firstExpr.Build(), matchExpr.Build())
}

func ISNULL(expression *Expression) string {
	return fmt.Sprintf("%s ISNULL", expression.Build())
}

func NOTNULL(expression *Expression) string {
	return fmt.Sprintf("%s NOTNULL", expression.Build())
}

func NOT_NULL(expression *Expression) string {
	return fmt.Sprintf("%s NOT NULL", expression.Build())
}

func IS_NOT_DISTINCT_FROM(firstExpr, secondExpr *Expression) string {
	return fmt.Sprintf("%s IS NOT DISTINCT FROM %s", firstExpr.Build(), secondExpr.Build())
}

func IS_DISTINCT_FROM(firstExpr, secondExpr *Expression) string {
	return fmt.Sprintf("%s IS DISTINCT FROM %s", firstExpr.Build(), secondExpr.Build())
}

func IS_NOT(firstExpr, secondExpr *Expression) string {
	return fmt.Sprintf("%s IS NOT %s", firstExpr.Build(), secondExpr.Build())
}

func IS(firstExpr, secondExpr *Expression) string {
	return fmt.Sprintf("%s IS %s", firstExpr.Build(), secondExpr.Build())
}

func NOT_BETWEEN(firstExpr, betweenExpr, andExpr *Expression) string {
	return fmt.Sprintf("%s NOT BETWEEN %s AND %s", firstExpr.Build(), betweenExpr.Build(), andExpr.Build())
}

func BETWEEN(firstExpr, betweenExpr, andExpr *Expression) string {
	return fmt.Sprintf("%s BETWEEN %s AND %s", firstExpr.Build(), betweenExpr.Build(), andExpr.Build())
}

// NOT_IN expression is not finished because of the complexity so be careful to use it
func NOT_IN(firstExpr *Expression, body string) string {
	return fmt.Sprintf("%s NOT IN %s", firstExpr.Build(), body)
}

// NOT_IN expression is not finished because of the complexity so be careful to use it
func IN(firstExpr *Expression, body string) string {
	return fmt.Sprintf("%s IN %s", firstExpr.Build(), body)
}

func NOT_EXISTS(sel *selectstmt.Select) string {
	return fmt.Sprintf("NOT EXISTS(%s)", sel.Build())
}

func EXISTS(sel *Select) string {
	return fmt.Sprintf("EXISTS(%s)", sel.Build())
}

func SELECTexpr(sel *Select) string {
	return fmt.Sprintf("(%s)", sel.Build())
}

func CASE(expr *Expression, whenExpr []*Expression, thenExpr []*Expression, elseExpr *Expression) string {
	var (
		whenThen = make([]string, len(whenExpr))
		elseStr  string
		exprStr  string
	)
	if len(whenExpr) != len(thenExpr) {
		panic("CASE statement requires an equal number of WHEN and THEN clauses")
	}

	if expr != nil {
		exprStr = fmt.Sprintf(" %s", expr.Build())
	}

	for i := range whenExpr {
		whenThen[i] = fmt.Sprintf("WHEN %s THEN %s", whenExpr[i].Build(), thenExpr[i].Build())
	}

	if elseExpr != nil {
		elseStr = fmt.Sprintf(" ELSE %s", elseExpr.Build())
	}

	return fmt.Sprintf("CASE%s %s%s END", exprStr, strings.Join(whenThen, " "), elseStr)
}

func RAISE(action ConflictClause, expr *Expression) string {
	if action == REPLACE || action == NO_CONFLICT {
		panic("RAISE statement cannot be used with REPLACE or empty string")
	}

	if action == IGNORE {
		return fmt.Sprintf("RAISE(%s)", action)
	}
	return fmt.Sprintf("RAISE(%s, %s)", action, expr.Build())
}
