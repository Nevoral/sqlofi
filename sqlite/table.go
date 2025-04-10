package sqlite

import (
	table "github.com/Nevoral/sqlofi/internal/sqlite/Table"
)

func CREATE_TABLE(model any, foreignTablePtrs ...any) *Table {
	return &Table{
		Table: table.NewTable(model, foreignTablePtrs),
	}
}

type Table struct {
	*table.Table
}

func (t *Table) Temporary() *Table {
	t.Table.Temporary()
	return t
}

func (t *Table) IfNotExists() *Table {
	t.Table.IfNotExists()
	return t
}

func (t *Table) Schema(schemaName string) *Table {
	t.Table.Schema(schemaName)
	return t
}

func (t *Table) Select(statement *Select) *Table {
	t.Table.Select(statement.Select)
	return t
}

func (t *Table) WithouRowID() *Table {
	t.Table.WithouRowID()
	return t
}

func (t *Table) Strict() *Table {
	t.Table.Strict()
	return t
}

// PrimaryKey creates a PRIMARY KEY constraint on the table.
// constaintName is the name of the constraint if "" it would be without named constraint,
// key is the PRIMARY KEY builder.
func (t *Table) PrimaryKey(constraintName string, key *PrimaryKey) *Table {
	t.Table.PrimaryKey(constraintName, key.TablePrimaryKey)
	return t
}

// Unique creates a UNIQUE constraint on the table.
// constaintName is the name of the constraint if "" it would be without named constraint,
// unique is the UNIQUE builder.
func (t *Table) Unique(constraintName string, unique *Unique) *Table {
	t.Table.Unique(constraintName, unique.TableUnique)
	return t
}

// Check creates a CHECK constraint on the table.
// constaintName is the name of the constraint if "" it would be without named constraint,
// expression is the expression to be included in the CHECK.
func (t *Table) Check(constraintName string, expression *Expression) *Table {
	t.Table.Check(constraintName, expression.Expression)
	return t
}

// ForeignKey creates a FOREIGN KEY constraint on the table.
// constaintName is the name of the constraint if "" it would be without named constraint,
// key is the FOREIGN KEY builders.
func (t *Table) ForeignKey(constraintName string, key *ForeignKey) *Table {
	t.Table.ForeignKey(constraintName, key.References)
	return t
}
