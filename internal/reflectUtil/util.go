package reflectutil

import "reflect"

// GetStructName returns the name of the struct referenced in the foreign key
func GetStructName(table any) string {
	var tableName string
	tableValue := reflect.ValueOf(table)

	if tableValue.Kind() == reflect.Ptr {
		tableValue = tableValue.Elem()
	}

	if tableValue.Kind() == reflect.Struct {
		tableName = tableValue.Type().Name()
	} else if str, ok := table.(string); ok {
		tableName = str
	}

	return tableName
}

// GetStructFieldsNames returns the names of all parameters/fields in the referenced struct
func GetStructFieldsNames(table any) []string {
	var fields []string
	tableValue := reflect.ValueOf(table)

	if tableValue.Kind() == reflect.Ptr {
		tableValue = tableValue.Elem()
	}

	if tableValue.Kind() == reflect.Struct {
		tableType := tableValue.Type()
		for i := range tableType.NumField() {
			field := tableType.Field(i)
			fields = append(fields, field.Name)
		}
	}

	return fields
}

func GetStructFields(table any) []reflect.StructField {
	var fields []reflect.StructField
	tableValue := reflect.ValueOf(table)

	if tableValue.Kind() == reflect.Ptr {
		tableValue = tableValue.Elem()
	}

	if tableValue.Kind() == reflect.Struct {
		tableType := tableValue.Type()
		for i := range tableType.NumField() {
			fields = append(fields, tableType.Field(i))
		}
	}

	return fields
}
