package sqlite

import (
	unique "github.com/Nevoral/sqlofi/internal/sqlite/Unique"
)

func UNIQUE(idxColumns ...*IndexedColumn) *Unique {
	return &Unique{
		TableUnique: unique.NewTableUnique(convertSliceOfIndexedColumn(idxColumns)),
	}
}

type Unique struct {
	*unique.TableUnique
}

func (u *Unique) OnConflict(conflict ConflictClause) *Unique {
	u.TableUnique.OnConflict(conflict.String())
	return u
}
