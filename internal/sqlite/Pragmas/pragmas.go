package pragmas

import "fmt"

func NewPragma(schemaName, name string) *Pragma {
	return &Pragma{
		schemaName: schemaName,
		name:       name,
		eqValue:    "",
		colValue:   "",
	}
}

type Pragma struct {
	schemaName string
	name       string
	eqValue    string
	colValue   string
}

func (p *Pragma) FuncType(value string) *Pragma {
	if p.eqValue != "" {
		panic("Error Pragma can be only one of the type func/value")
	}
	p.colValue = value
	return p
}

func (p *Pragma) ValueType(value string) *Pragma {
	if p.colValue != "" {
		panic("Error Pragma can be only one of the type func/value")
	}
	p.eqValue = value
	return p
}

func (p *Pragma) Build() string {
	var (
		sch   string
		value string
	)
	if p.schemaName != "" {
		sch = p.schemaName + "."
	}

	if p.eqValue != "" {
		value = " = " + p.eqValue
	}
	if p.colValue != "" {
		value = " (" + p.eqValue + ")"
	}
	return fmt.Sprintf("PRAGMA %s%s%s", sch, p.name, value)
}
