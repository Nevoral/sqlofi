package defaultConstr

import (
	"fmt"
	"strings"
	"unicode"

	expr "github.com/Nevoral/sqlofi/internal/sqlite/Expression"
	types "github.com/Nevoral/sqlofi/internal/sqlite/Types"
)

func NewDefaultValue[T expr.Expression | types.LiteralValue | types.SignedNumber](value T) string {
	if val, ok := any(value).(expr.Expression); ok {
		return fmt.Sprintf("(%s)", val.Build())
	} else if val, ok := any(value).(types.LiteralValue); ok {
		return fmt.Sprintf("%s", val.String())
	} else if val, ok := any(value).(types.SignedNumber); ok {
		return fmt.Sprintf("%s", val.String())
	} else {
		panic(fmt.Sprintf("unsupported default value type: %T", value))
	}
}

// ParseDefault parses a default constraint value from a string
// Expected format: "NULL", "TRUE", "FALSE", number, 'string', CURRENT_TIME, etc.
func ParseDefault(content string) string {
	content = strings.TrimSpace(content)

	// Handle empty case
	if content == "" {
		return ""
	}

	// Remove DEFAULT keyword if present
	if strings.HasPrefix(strings.ToUpper(content), "DEFAULT") {
		content = strings.TrimSpace(content[len("DEFAULT"):])
	}

	defaultVal := parseDefaultValue(content)

	return fmt.Sprintf("DEFAULT %s", defaultVal)
}

// parseDefaultValue determines the type and value of the default
func parseDefaultValue(content string) string {
	content = strings.TrimSpace(content)

	// Handle NULL
	if strings.ToUpper(content) == "NULL" {
		return NewDefaultValue(types.NULL_VALUE)
	}

	// Handle boolean literals
	if strings.ToUpper(content) == "TRUE" {
		return NewDefaultValue(types.TRUE_VALUE)
	}
	if strings.ToUpper(content) == "FALSE" {
		return NewDefaultValue(types.FALSE_VALUE)
	}

	// Handle CURRENT_TIME, CURRENT_DATE, CURRENT_TIMESTAMP
	if strings.HasPrefix(strings.ToUpper(content), "CURRENT_TIME") {
		return NewDefaultValue(types.CURRENT_TIME_VALUE)
	}
	if strings.HasPrefix(strings.ToUpper(content), "CURRENT_DATE") {
		return NewDefaultValue(types.CURRENT_DATE_VALUE)
	}
	if strings.HasPrefix(strings.ToUpper(content), "CURRENT_TIMESTAMP") {
		return NewDefaultValue(types.CURRENT_TIMESTAMP_VALUE)
	}

	// Handle string literals (enclosed in quotes)
	if (strings.HasPrefix(content, "'") && strings.HasSuffix(content, "'")) ||
		(strings.HasPrefix(content, "\"") && strings.HasSuffix(content, "\"")) {
		// Remove surrounding quotes
		value := content[1 : len(content)-1]
		return NewDefaultValue(types.LiteralValue(value))
	}

	// Handle expressions (enclosed in parentheses)
	if strings.HasPrefix(content, "(") && strings.HasSuffix(content, ")") {
		// Remove surrounding parentheses
		e := content[1 : len(content)-1]
		return NewDefaultValue(*expr.NewExpression(e))
	}

	// Handle numbers (including signed numbers like +10, -5)
	isNumber := true

	for i, c := range content {
		if i == 0 && (c == '+' || c == '-') {
			continue
		}

		if c == '.' {
			continue
		}

		if !unicode.IsDigit(c) {
			isNumber = false
			break
		}
	}

	if isNumber {
		return NewDefaultValue(types.NewSignedNumber(content))
	}

	// If we can't determine the type, just return the content as is
	// This could be a function call or other expression not in parentheses
	return NewDefaultValue(*expr.NewExpression(content))
}
