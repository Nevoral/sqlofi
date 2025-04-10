package idxcol

import (
	"fmt"

	expr "github.com/Nevoral/sqlofi/internal/sqlite/Expression"
	sortorder "github.com/Nevoral/sqlofi/internal/sqlite/SortOrder"
	"github.com/Nevoral/sqlofi/internal/utils"
)

func NewIndexedColumnNames(name string) *IndexedColumn {
	return &IndexedColumn{
		name: name,
	}
}

func NewIndexedColumnExpresions(expresion *expr.Expression) *IndexedColumn {
	return &IndexedColumn{
		expression: expresion,
	}
}

type IndexedColumn struct {
	name       string
	collate    string
	sortOrder  sortorder.SortOrder
	expression *expr.Expression
}

func (i *IndexedColumn) Collate(name string) *IndexedColumn {
	i.collate = name
	return i
}

func (i *IndexedColumn) ASC() *IndexedColumn {
	if i.sortOrder == "" {
		i.sortOrder = sortorder.ASC
	}
	return i
}

func (i *IndexedColumn) DESC() *IndexedColumn {
	if i.sortOrder == "" {
		i.sortOrder = sortorder.DESC
	}
	return i
}

func (i *IndexedColumn) Build() string {
	var (
		start  string
		middle string
		end    string
	)

	if i.name == "" {
		start = i.expression.Build()
	} else {
		start = utils.ToSnakeCase(i.name)
	}

	if i.collate != "" {
		middle = fmt.Sprintf(" COLLATE %s", i.collate)
	}

	if i.sortOrder != "" {
		end = fmt.Sprintf(" %s", i.sortOrder.String())
	}
	return fmt.Sprintf("%s%s%s", start, middle, end)
}
