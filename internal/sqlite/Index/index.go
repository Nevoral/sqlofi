package index

import (
	"fmt"

	reflectutil "github.com/Nevoral/sqlofi/internal/reflectUtil"
	expr "github.com/Nevoral/sqlofi/internal/sqlite/Expression"
	idxcol "github.com/Nevoral/sqlofi/internal/sqlite/IndexedColumn"
)

func NewIndex(table any, indexName string, indexCols []*idxcol.IndexedColumn) *Index {
	return &Index{
		name:    indexName,
		table:   table,
		columns: indexCols,
	}
}

type Index struct {
	unique      bool
	ifNotExists bool
	schemaName  string
	name        string
	table       any
	columns     []*idxcol.IndexedColumn
	where       *expr.Expression
}

func (i *Index) Unique() *Index {
	i.unique = true
	return i
}

func (i *Index) IfNotExists() *Index {
	i.ifNotExists = true
	return i
}

func (i *Index) Schema(schemaName string) *Index {
	i.schemaName = schemaName
	return i
}

func (i *Index) Where(expression *expr.Expression) *Index {
	i.where = expression
	return i
}

func (i *Index) Build() string {
	var (
		uniq   string
		ifnot  string
		schema string
		col    string
		where  string
	)
	if i.unique {
		uniq = "UNIQUE "
	}
	if i.ifNotExists {
		ifnot = "IF NOT EXISTS "
	}
	if i.schemaName != "" {
		schema = fmt.Sprintf("%s.", i.schemaName)
	}
	for _, column := range i.columns {
		col += fmt.Sprintf("%s, ", column.Build())
	}
	if i.where != nil {
		where = fmt.Sprintf(" WHERE %s", i.where.Build())
	}
	return fmt.Sprintf("CREATE %sINDEX %s%s%s ON %s (%s)%s", uniq, ifnot, schema, i.name, reflectutil.GetStructName(i.table), col, where)
}
