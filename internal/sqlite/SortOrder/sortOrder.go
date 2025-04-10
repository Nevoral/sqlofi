package sortorder

import "fmt"

func NewSortOrder(value string) (SortOrder, error) {
	switch value {
	case "":
		return UNSORTED, nil
	case "ASC":
		return ASC, nil
	case "DESC":
		return DESC, nil
	default:
		return UNSORTED, fmt.Errorf("invalid sort order: %s", value)
	}
}

type SortOrder string

const (
	UNSORTED SortOrder = ""
	ASC      SortOrder = "ASC"
	DESC     SortOrder = "DESC"
)

func (so SortOrder) String() string {
	return string(so)
}
