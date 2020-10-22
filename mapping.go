package mappert

import (
	"fmt"
	"reflect"
)

func followPointer(in reflect.Value) reflect.Value {

	found := in

	for found.Kind() == reflect.Ptr {
		found = found.Elem()
	}

	return found
}

func followPointerUntil(in reflect.Value, untilKind reflect.Kind) (reflect.Value, error) {

	in = followPointer(in)

	if kind := in.Kind(); kind != untilKind {
		return reflect.ValueOf(nil), fmt.Errorf("expect a %s kind, not %s", untilKind, kind)
	}

	return in, nil
}

func mapFieldValue(in reflect.Value) interface{} {

	in = followPointer(in)

	// let's see how we map
	//
	switch in.Kind() {
	case reflect.Slice:
		return slice2SliceOfMaps(in)
	case reflect.Struct:
		return struct2Map(in)
	default:
		return in.Interface()
	}

}

func slice2SliceOfMaps(in reflect.Value) []interface{} {

	l := in.Len()

	result := make([]interface{}, l, l)

	for index := 0; index < l; index++ {
		result[index] = mapFieldValue(in.Index(index))
	}

	return result
}

func struct2Map(in reflect.Value) map[string]interface{} {
	structType := in.Type()

	result := make(map[string]interface{}, 0)
	for fieldIndex := 0; fieldIndex < structType.NumField(); fieldIndex++ {
		fieldName := structType.Field(fieldIndex).Name
		fieldValue := in.Field(fieldIndex)

		result[fieldName] = mapFieldValue(fieldValue)
	}

	return result

}

// Struct2Map converts a pointer to a struct to a map from field names to field values
//
func Struct2Map(in interface{}) (map[string]interface{}, error) {

	valueIn, err := followPointerUntil(reflect.ValueOf(in), reflect.Struct)
	if err != nil {
		return nil, err
	}

	result := struct2Map(valueIn)

	return result, nil
}

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
