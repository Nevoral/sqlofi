package sqlite

import (
	"fmt"

	foreignkey "github.com/Nevoral/sqlofi/internal/sqlite/ForeignKey"
)

type RowAction string

const (
	CASCADE     RowAction = "CASCADE"
	SET_NULL    RowAction = "SET NULL"
	SET_DEFAULT RowAction = "SET DEFAULT"
	RESTRICT    RowAction = "RESTRICT"
	NO_ACTION   RowAction = "NO ACTION"
)

func (r RowAction) String() string {
	return string(r)
}

type DeferrableAction string

const (
	NO_DEFERRABLE_ACTION DeferrableAction = ""
	INITIALLY_DEFERRED   DeferrableAction = "INITIALLY DEFERRED"
	INITIALLY_IMMEDIATE  DeferrableAction = "INITIALLY IMMEDIATE"
)

func (d DeferrableAction) String() string {
	return string(d)
}

func FOREIGN_KEY(foreignTablePtr any, columns ...string) *ForeignKey {
	if foreignTablePtr == nil {
		panic(fmt.Errorf("Error no Table provided"))
	}
	return &ForeignKey{
		References: foreignkey.NewTableForeignTable(foreignTablePtr, columns),
	}
}

type ForeignKey struct {
	*foreignkey.References
}

func (f *ForeignKey) OnDelete(action RowAction) *ForeignKey {
	f.References.OnDelete(action.String())
	return f
}

func (f *ForeignKey) OnUpdate(action RowAction) *ForeignKey {
	f.References.OnUpdate(action.String())
	return f
}

func (f *ForeignKey) Match(name string) *ForeignKey {
	f.References.Match(name)
	return f
}

func (f *ForeignKey) Deferrable(action DeferrableAction) *ForeignKey {
	f.References.Deferrable(action.String())
	return f
}

func (f *ForeignKey) NotDeferrable(action DeferrableAction) *ForeignKey {
	f.References.NotDeferrable(action.String())
	return f
}

func (f *ForeignKey) ForeighColumns(columns ...string) *ForeignKey {
	f.References.ForeighColumns(columns)
	return f
}
