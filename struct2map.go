package mappert

import (
	"reflect"
)

func mapFieldValue(in reflect.Value, config *MapConfiguration) interface{} {

	in = followPointer(in)

	// let's see how we map
	//
	switch in.Kind() {
	case reflect.Slice:
		return slice2SliceOfMaps(in, config)
	case reflect.Struct:
		return struct2Map(in, config)
	default:
		return in.Interface()
	}

}

func slice2SliceOfMaps(in reflect.Value, config *MapConfiguration) []interface{} {

	l := in.Len()

	result := make([]interface{}, l, l)

	for index := 0; index < l; index++ {
		result[index] = mapFieldValue(in.Index(index), config)
	}

	return result
}

func struct2Map(in reflect.Value, config *MapConfiguration) map[string]interface{} {
	structType := in.Type()

	result := make(map[string]interface{}, 0)
	for fieldIndex := 0; fieldIndex < structType.NumField(); fieldIndex++ {
		fieldName := structType.Field(fieldIndex).Name
		if config.ShouldIgnore(fieldName) {
			continue
		}

		result[config.Name(fieldName)] = mapFieldValue(in.Field(fieldIndex), config)
	}

	return result

}

// Struct2Map converts a pointer to a struct to a map from field names to field values
//
func Struct2Map(in interface{}, config ...*MapConfiguration) (map[string]interface{}, error) {

	valueIn, err := followPointerUntil(reflect.ValueOf(in), reflect.Struct)
	if err != nil {
		return nil, err
	}

	result := struct2Map(valueIn, combineConfiguration(config...))

	return result, nil
}
