package sqlite

import (
	idxcol "github.com/Nevoral/sqlofi/internal/sqlite/IndexedColumn"
)

func NewIndexedColumn[T ~string | *Expression](value T) *IndexedColumn {
	if str, ok := any(value).(string); ok {
		return &IndexedColumn{
			IndexedColumn: idxcol.NewIndexedColumnNames(str),
		}
	}
	if expr, ok := any(value).(*Expression); ok {
		return &IndexedColumn{
			IndexedColumn: idxcol.NewIndexedColumnExpresions(expr.Expression),
		}
	}
	return nil
}

type IndexedColumn struct {
	*idxcol.IndexedColumn
}

func convertSliceOfIndexedColumn(columns []*IndexedColumn) []*idxcol.IndexedColumn {
	var result []*idxcol.IndexedColumn
	for _, column := range columns {
		result = append(result, column.IndexedColumn)
	}
	return result
}
