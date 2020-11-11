package mappert

import (
	"fmt"
	"reflect"
)

// Map2Struct assigns all fields of a given struct with corresponding map data
//
func Map2Struct(in map[string]interface{}, out interface{}, configs ...*MapConfiguration) {

	config := combineConfiguration(configs...)

	for key, value := range in {
		name := config.Name(key)
		if config.ShouldIgnore(name) {
			continue
		}
		setFieldValue(reflect.ValueOf(out).Elem(), name, reflect.ValueOf(value), config)
	}
}

func createSliceOfStructsFromSliceOfMaps(field reflect.Value, fieldValue reflect.Value) reflect.Value {

	if fieldValue.Kind() != reflect.Slice {
		panic(fmt.Sprintln("can not assign", fieldValue, "to", field))
	}

	size := fieldValue.Len()
	created := reflect.MakeSlice(field.Type(), size, size)
	for index := 0; index < size; index++ {
		Map2Struct(fieldValue.Index(index).Interface().(map[string]interface{}), created.Index(index).Addr().Interface())
	}

	return created
}

func setFieldValue(out reflect.Value, fieldName string, fieldValue reflect.Value, config *MapConfiguration) {

	field := out.FieldByName(fieldName)
	if !field.IsValid() {
		return
	}

	switch field.Kind() {
	case reflect.Slice:
		if sliceKind(field) == reflect.Struct {
			field.Set(createSliceOfStructsFromSliceOfMaps(field, fieldValue))
		} else {
			field.Set(fieldValue)
		}

	case reflect.Struct:
		innerValue := reflect.New(field.Type())

		Map2Struct(fieldValue.Interface().(map[string]interface{}), innerValue.Interface(), config)

		field.Set(innerValue.Elem())
	default:
		field.Set(fieldValue)
	}

}
