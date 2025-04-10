package sqlite

import index "github.com/Nevoral/sqlofi/internal/sqlite/Index"

func CREATE_INDEX(table any, indexName string, indexCols ...*IndexedColumn) *Index {
	return &Index{
		Index: index.NewIndex(table, indexName, convertSliceOfIndexedColumn(indexCols)),
	}
}

type Index struct {
	*index.Index
}

func (i *Index) Unique() *Index {
	i.Index.Unique()
	return i
}

func (i *Index) IfNotExists() *Index {
	i.Index.IfNotExists()
	return i
}

func (i *Index) Schema(schemaName string) *Index {
	i.Index.Schema(schemaName)
	return i
}

func (i *Index) Where(expression *Expression) *Index {
	i.Index.Where(expression.Expression)
	return i
}
