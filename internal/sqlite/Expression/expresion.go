package expr

func NewExpression(expression string) *Expression {
	return &Expression{
		expression: expression,
	}
}

// Expression represents an SQL expression
type Expression struct {
	expression string
}

// Build returns the string representation of the expression
func (e *Expression) Build() string {
	return e.expression
}
