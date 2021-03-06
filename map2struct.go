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
		if config.ShouldIgnore(key) {
			continue
		}
		setFieldValue(reflect.ValueOf(out).Elem(), config.Name(key), reflect.ValueOf(value), config)
	}
}

func createSliceOfStructsFromSliceOfMaps(field reflect.Value, fieldValue reflect.Value, config *MapConfiguration) reflect.Value {

	if fieldValue.Kind() != reflect.Slice {
		panic(fmt.Sprintln("can not assign", fieldValue, "to", field))
	}

	size := fieldValue.Len()
	created := reflect.MakeSlice(field.Type(), size, size)
	for index := 0; index < size; index++ {
		Map2Struct(fieldValue.Index(index).Interface().(map[string]interface{}), created.Index(index).Addr().Interface(), config)
	}

	return created
}

func convertAndSetFieldValue(field reflect.Value, fieldName string, fieldValue reflect.Value, config *MapConfiguration) {

	if field.Type() != fieldValue.Type() {
		converted, err := config.ConvertValue(fieldValue.Interface(), field.Type())

		if err != nil {
			panic(fmt.Sprintln("field", fieldName, "expects", field.Type(), "/", field.Kind(), "not", fieldValue, "which is", fieldValue.Type(), "/", fieldValue.Kind(), ":", err))
		}

		convertedValue := reflect.ValueOf(converted)

		if field.Type() != convertedValue.Type() {
			panic(fmt.Sprintln("field", fieldName, "expects", field.Type(), "/", field.Kind(), "not", convertedValue, "which is", convertedValue.Type(), "/", convertedValue.Kind(), "(converted from", fieldValue.Type(), "/", fieldValue.Kind(), ")"))
		}
		field.Set(convertedValue)
	} else {
		field.Set(fieldValue)
	}
}

func setFieldValue(out reflect.Value, fieldName string, fieldValue reflect.Value, config *MapConfiguration) {

	field := out.FieldByName(fieldName)
	if !field.IsValid() {
		return
	}

	switch field.Kind() {
	case reflect.Slice:
		if sliceKind(field) == reflect.Struct {
			field.Set(createSliceOfStructsFromSliceOfMaps(field, fieldValue, config))
		} else {
			field.Set(fieldValue)
		}

	case reflect.Struct:

		if fieldValueAsMap, couldCastToMap := fieldValue.Interface().(map[string]interface{}); couldCastToMap {
			innerValue := reflect.New(field.Type())
			Map2Struct(fieldValueAsMap, innerValue.Interface(), config)
			field.Set(innerValue.Elem())
		} else {
			convertAndSetFieldValue(field, fieldName, fieldValue, config)
		}
	default:
		convertAndSetFieldValue(field, fieldName, fieldValue, config)
	}

}
