package utils

import "strings"

func ToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func Join(elems []string, sep string) (out string) {
	if len(elems) == 0 {
		return
	}

	for idx, elem := range elems {
		if idx != 0 {
			out += ", "
		}
		out += ToSnakeCase(elem)
	}
	return
}
