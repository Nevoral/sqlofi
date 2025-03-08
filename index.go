package sqlofi

import (
	"fmt"
	"strings"
)

func parseIndex(unique bool, value string) *IndexOptions {
	name, where, _ := strings.Cut(value, ",WHERE ")
	return &IndexOptions{
		Name:   name,
		Unique: unique,
		Where:  where,
	}
}

type IndexOptions struct {
	Name   string // Name of the index
	Unique bool   // Whether the index is unique
	Where  string // Conditional index (partial index)
	Column string
}

func (i *IndexOptions) String(table string) string {
	var (
		un = " UNIQUE"
		wr = " WHERE " + i.Where
	)
	if !i.Unique {
		un = ""
	}

	if i.Where == "" {
		wr = ""
	}
	return fmt.Sprintf("CREATE %sINDEX %s ON %s (%s)%s;", un, i.Name, table, i.Column, wr)
}
