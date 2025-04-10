package unique

import (
	"fmt"
	"strings"

	idxcol "github.com/Nevoral/sqlofi/internal/sqlite/IndexedColumn"
)

func NewColumnUnique(conflict string) string {
	var con string
	if conflict != "" {
		con = " " + conflict
	}
	return fmt.Sprintf("UNIQUE%s", con)
}

func NewTableUnique(idxColumns []*idxcol.IndexedColumn) *TableUnique {
	return &TableUnique{
		indexedColumn: idxColumns,
	}
}

type TableUnique struct {
	indexedColumn []*idxcol.IndexedColumn
	conflict      string
}

func (u *TableUnique) OnConflict(conflict string) *TableUnique {
	u.conflict = conflict
	return u
}

func (u *TableUnique) Build() string {
	var col []string
	for _, val := range u.indexedColumn {
		col = append(col, val.Build())
	}

	return fmt.Sprintf("UNIQUE (%s) %s", strings.Join(col, ", "), u.conflict)
}
