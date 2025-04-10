package check

import (
	"fmt"

	expr "github.com/Nevoral/sqlofi/internal/sqlite/Expression"
)

func NewCheck(ex *expr.Expression) *Check {
	return &Check{ex}
}

type Check struct {
	*expr.Expression
}

func (c *Check) Build() string {
	return fmt.Sprintf("CHECK (%s)", c.Expression.Build())
}
