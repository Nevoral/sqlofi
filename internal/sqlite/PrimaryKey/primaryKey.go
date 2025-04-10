package primarykey

import (
	"fmt"
	"strings"

	idxcol "github.com/Nevoral/sqlofi/internal/sqlite/IndexedColumn"
	sortorder "github.com/Nevoral/sqlofi/internal/sqlite/SortOrder"
)

func NewTablePrimaryKey(idxColumns []*idxcol.IndexedColumn) *TablePrimaryKey {
	return &TablePrimaryKey{
		indexedColumn: idxColumns,
	}
}

type TablePrimaryKey struct {
	indexedColumn []*idxcol.IndexedColumn
	conflict      string
}

func (t *TablePrimaryKey) OnConflict(conflict string) *TablePrimaryKey {
	t.conflict = conflict
	return t
}

func (t *TablePrimaryKey) Build() string {
	var (
		col      []string
		conflict string
	)

	for _, val := range t.indexedColumn {
		col = append(col, val.Build())
	}

	if t.conflict != "" {
		conflict = fmt.Sprintf(" %s", t.conflict)
	}

	return fmt.Sprintf("PRIMARY KEY (%s)%s", strings.Join(col, ", "), conflict)
}

func NewColumnPrimaryKey(sortOrder sortorder.SortOrder) *ColumnPrimaryKey {
	return &ColumnPrimaryKey{
		sortOrder: sortOrder,
	}
}

type ColumnPrimaryKey struct {
	sortOrder     sortorder.SortOrder
	conflict      string
	autoincrement bool
}

func (c *ColumnPrimaryKey) Conflict(conflict string) *ColumnPrimaryKey {
	c.conflict = conflict
	return c
}

func (c *ColumnPrimaryKey) Autoincrement(value bool) *ColumnPrimaryKey {
	c.autoincrement = value
	return c
}

func (c *ColumnPrimaryKey) Build() string {
	var (
		autoincrement string
		sort          string
		conflict      string
	)
	if c.sortOrder != sortorder.UNSORTED {
		sort = fmt.Sprintf(" %s", c.sortOrder.String())
	}

	if c.conflict != "" {
		conflict = fmt.Sprintf(" %s", c.conflict)
	}

	if c.autoincrement {
		autoincrement = " AUTOINCREMENT"
	}

	return fmt.Sprintf("PRIMARY KEY%s%s%s", sort, conflict, autoincrement)
}
