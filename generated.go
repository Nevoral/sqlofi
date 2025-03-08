package sqlofi

import (
	"fmt"
	"strings"
)

func parseGenerated(value string) *GeneratedOption {
	logic, found := strings.CutSuffix(value, ",STORED")
	return &GeneratedOption{
		computeLogic: logic,
		Stored:       found,
	}
}

type GeneratedOption struct {
	computeLogic string
	Stored       bool
}

func (g *GeneratedOption) String() string {
	if g.Stored {
		return fmt.Sprintf("GENERATED ALWAYS AS %s STORED", g.computeLogic)
	}
	return fmt.Sprintf("GENERATED ALWAYS AS %s VIRTUAL", g.computeLogic)
}
