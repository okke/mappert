package mappert

import "reflect"

func setFieldValue(out reflect.Value, fieldName string, fieldValue reflect.Value) {

	field := out.FieldByName(fieldName)

	switch field.Kind() {
	case reflect.Struct:
		innerValue := reflect.New(field.Type())

		Map2Struct(fieldValue.Interface().(map[string]interface{}), innerValue.Interface())

		field.Set(innerValue.Elem())
	default:
		field.Set(fieldValue)
	}

}

// Map2Struct assigns all fields of a given struct with corresponding map data
//
func Map2Struct(in map[string]interface{}, out interface{}) {

	for key, value := range in {
		setFieldValue(reflect.ValueOf(out).Elem(), key, reflect.ValueOf(value))
	}
}
