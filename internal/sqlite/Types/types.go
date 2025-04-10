package types

import (
	"database/sql"
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	reflectutil "github.com/Nevoral/sqlofi/internal/reflectUtil"
)

// getSQLiteType returns the appropriate SQLite type for a Go type
func GetSQLiteType(t reflect.Type) (SQLiteType, bool) {
	var (
		sqlType    = TEXT
		foreignKey bool
	)

	// Handle pointer types
	if t.Kind() == reflect.Ptr {
		elemType, isForeignKey := GetSQLiteType(t.Elem())
		return elemType, isForeignKey || t.Elem().Kind() == reflect.Struct
	}

	// Handle sql.Null* types
	switch t {
	case reflect.TypeOf(sql.NullBool{}):
		return "INTEGER", foreignKey
	case reflect.TypeOf(sql.NullInt16{}):
		return "INTEGER", foreignKey
	case reflect.TypeOf(sql.NullInt32{}):
		return "INTEGER", foreignKey
	case reflect.TypeOf(sql.NullInt64{}):
		return "INTEGER", foreignKey
	case reflect.TypeOf(sql.NullFloat64{}):
		return "REAL", foreignKey
	case reflect.TypeOf(sql.NullString{}):
		return "TEXT", foreignKey
	case reflect.TypeOf(sql.NullTime{}):
		return "TEXT", foreignKey
	case reflect.TypeOf(sql.NullByte{}):
		return "BLOB", foreignKey
	}

	// Handle basic Go types
	switch t.Kind() {
	case reflect.Bool:
		return "INTEGER", foreignKey
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "INTEGER", foreignKey
	case reflect.Float32, reflect.Float64:
		return "REAL", foreignKey
	case reflect.String:
		return "TEXT", foreignKey
	case reflect.Struct:
	case reflect.Interface:
	case reflect.Map:
	case reflect.Array:

	}

	// Handle special types
	switch t {
	case reflect.TypeOf(time.Time{}):
		return "TEXT", foreignKey
	case reflect.TypeOf([]byte{}):
		return "BLOB", foreignKey
	}

	return sqlType, foreignKey
}

func NewSQLiteType(value string) SQLiteType {
	value = strings.TrimSpace(value)
	switch value {
	case "NULL":
		return NULL
	case "TEXT":
		return TEXT
	case "INTEGER":
		return INTEGER
	case "REAL":
		return REAL
	case "BLOB":
		return BLOB
	default:
		return NULL
	}
}

type SQLiteType string

const (
	NULL    SQLiteType = "NULL"
	TEXT    SQLiteType = "TEXT"
	INTEGER SQLiteType = "INTEGER"
	REAL    SQLiteType = "REAL"
	BLOB    SQLiteType = "BLOB"
)

func (t SQLiteType) String() string {
	if t == "" {
		return "NULL"
	}
	return string(t)
}

type LiteralValue string

const (
	NULL_VALUE              LiteralValue = "NULL"
	TRUE_VALUE              LiteralValue = "TRUE"
	FALSE_VALUE             LiteralValue = "FALSE"
	CURRENT_TIME_VALUE      LiteralValue = "CURRENT_TIME"
	CURRENT_DATE_VALUE      LiteralValue = "CURRENT_DATE"
	CURRENT_TIMESTAMP_VALUE LiteralValue = "CURRENT_TIMESTAMP"
)

func (l LiteralValue) String() string {
	if l == "" {
		return "NULL"
	}
	return string(l)
}

func NewSignedNumber[T string | int | uint | float32 | float64 | []byte | int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64](value T) SignedNumber {
	var result string

	switch v := any(value).(type) {
	case string:
		// For string, just return it directly
		return SignedNumber(v)
	case []byte:
		// Convert byte slice to string
		return SignedNumber(string(v))
	case int, int8, int16, int32, int64:
		// For signed integer types
		var val int64
		switch typedVal := any(value).(type) {
		case int:
			val = int64(typedVal)
		case int8:
			val = int64(typedVal)
		case int16:
			val = int64(typedVal)
		case int32:
			val = int64(typedVal)
		case int64:
			val = typedVal
		}

		if val >= 0 {
			result = "+ " + strconv.FormatInt(val, 10)
		} else {
			result = "- " + strconv.FormatInt(-val, 10)
		}
	case uint, uint8, uint16, uint32, uint64:
		// For unsigned integer types (always positive)
		var val uint64
		switch typedVal := any(value).(type) {
		case uint:
			val = uint64(typedVal)
		case uint8:
			val = uint64(typedVal)
		case uint16:
			val = uint64(typedVal)
		case uint32:
			val = uint64(typedVal)
		case uint64:
			val = typedVal
		}

		result = "+ " + strconv.FormatUint(val, 10)
	case float32, float64:
		// For floating point types
		var val float64
		switch typedVal := any(value).(type) {
		case float32:
			val = float64(typedVal)
		case float64:
			val = typedVal
		}

		if val >= 0 {
			result = "+ " + strconv.FormatFloat(val, 'f', -1, 64)
		} else {
			result = "- " + strconv.FormatFloat(-val, 'f', -1, 64)
		}
	}

	return SignedNumber(result)
}

type SignedNumber string

func (s SignedNumber) String() string {
	return string(s)
}

func NewDbPath(schema string, table any, column string) *DbPath {
	if table != nil {
		if !slices.Contains(reflectutil.GetStructFieldsNames(table), column) {
			panic("Error column name isn't present in the table")
		}
	}

	return &DbPath{
		schema: schema,
		table:  table,
		column: column,
	}
}

type DbPath struct {
	schema string
	table  any
	column string
}

func (d *DbPath) StringColumn() string {
	if d.schema == "" {
		if d.table == nil {
			return fmt.Sprintf("%s", d.column)
		}
		return fmt.Sprintf("%s.%s", d.table, d.column)
	}
	return fmt.Sprintf("%s.%s.%s", d.schema, d.table, d.column)
}

func (d *DbPath) StringTable() string {
	if d.schema == "" {
		return fmt.Sprintf("%s", d.table)
	}
	return fmt.Sprintf("%s.%s", d.schema, d.table)
}

func (d *DbPath) StringSchema() string {
	return fmt.Sprintf("%s", d.schema)
}

type (
	BindingParameter string
	BindingType      int
)

const (
	AUTOINCREMENTED BindingType = iota
	INDEXED
	COLON_NAMED
	AT_NAMED
	DOLAR_NAMED
)

func NewBindingParameter(b BindingType, value string) BindingParameter {
	switch b {
	case INDEXED:
		return BindingParameter("?" + value)
	case COLON_NAMED:
		return BindingParameter(":" + value)
	case AT_NAMED:
		return BindingParameter("@" + value)
	case DOLAR_NAMED:
		return BindingParameter("$" + value)
	default:
		return BindingParameter("?")
	}
}

func (b BindingParameter) String() string {
	return string(b)
}
