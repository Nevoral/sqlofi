package notnull

func NewNotNull(conflict string) string {
	if conflict != "" {
		return "NOT NULL " + conflict
	}
	return "NOT NULL"
}
