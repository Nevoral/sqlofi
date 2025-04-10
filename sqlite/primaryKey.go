package sqlite

import (
	primarykey "github.com/Nevoral/sqlofi/internal/sqlite/PrimaryKey"
)

func PRIMARY_KEY(idxColumns ...*IndexedColumn) *PrimaryKey {
	return &PrimaryKey{
		TablePrimaryKey: primarykey.NewTablePrimaryKey(convertSliceOfIndexedColumn(idxColumns)),
	}
}

type PrimaryKey struct {
	*primarykey.TablePrimaryKey
}

func (p *PrimaryKey) OnConflict(conflict ConflictClause) *PrimaryKey {
	p.TablePrimaryKey.OnConflict(conflict.String())
	return p
}
