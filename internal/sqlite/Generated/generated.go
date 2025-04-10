package generated

import (
	"fmt"

	expr "github.com/Nevoral/sqlofi/internal/sqlite/Expression"
)

type StorageType string

const (
	NO_STORAGE StorageType = ""
	STORED     StorageType = "STORED"
	VIRTUAL    StorageType = "VIRTUAL"
)

func (s StorageType) String() string {
	return string(s)
}

func NewGenerated(always bool, expr *expr.Expression, storage StorageType) string {
	var (
		alw    string
		stored string
	)

	if always {
		alw = "GENERATED ALWAYS "
	}

	if storage != NO_STORAGE {
		stored = storage.String()
	}

	return fmt.Sprintf("%sAS (%s)%s", alw, expr.Build(), stored)
}
