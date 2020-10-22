package mappert

import "reflect"

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
