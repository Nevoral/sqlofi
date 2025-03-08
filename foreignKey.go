package sqlofi

import (
	"fmt"
	"strings"
)

type RowAction int

const (
	Cascade RowAction = iota
	SetNul
	SetDefault
	Restrict
)

func (r RowAction) String() string {
	return []string{"CASCADE", "SET NULL", "SET DEFAULT", "RESTRICT"}[r]
}

type Action func(RowAction) string

var (
	OnDelete = func(r RowAction) string {
		return fmt.Sprintf("ON DELETE %s", r.String())
	}
	OnUpdate = func(r RowAction) string {
		return fmt.Sprintf("ON UPDATE %s", r.String())
	}
	None = func(r RowAction) string {
		return ""
	}
)

func parseForeignKey(value string) *ForeignKey {
	reference, action, found := strings.Cut(value, ",")
	if found {
		//TODO: cheking corectness of action
	}
	return &ForeignKey{
		References: reference,
		Action:     action,
	}
}

type ForeignKey struct {
	Column     string
	References string
	Action     string
}

func (f *ForeignKey) String() string {
	if f.Action == "" {
		return fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s %s", f.Column, f.References, f.Action)
	}
	return fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s", f.Column, f.References)
}
