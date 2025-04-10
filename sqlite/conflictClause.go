package sqlite

type ConflictClause string

const (
	NO_CONFLICT ConflictClause = ""
	ROLLBACK    ConflictClause = "ROLLBACK"
	ABORT       ConflictClause = "ABORT"
	FAIL        ConflictClause = "FAIL"
	IGNORE      ConflictClause = "IGNORE"
	REPLACE     ConflictClause = "REPLACE"
)

func (c ConflictClause) String() string {
	if c == NO_CONFLICT {
		return ""
	}
	return "ON CONFLICT " + string(c)
}
