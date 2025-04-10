package collate

import "fmt"

func NewCollate(name string) string {
	return fmt.Sprintf("COLLATE %s", name)
}
