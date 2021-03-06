package mappert

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"

	"github.com/relvacode/iso8601"
)

// ConvertorFunc is a function that maps a value to a value of a different type
//
type ConvertorFunc func(in interface{}) (interface{}, error)

// Convertor is a registry of convertor functions mapped by type
//
type Convertor interface {
	RegisterConvertor(from, to reflect.Type, convertor ConvertorFunc) Convertor

	ConvertIn(in, out interface{}) error
	ConvertTo(in interface{}, to reflect.Type) (interface{}, error)
}

type convertorMapping struct {
	mapping map[reflect.Type]map[reflect.Type]ConvertorFunc
}

func (convertorMapping *convertorMapping) RegisterConvertor(from, to reflect.Type, convertor ConvertorFunc) Convertor {

	if convertorMapping.mapping == nil {
		convertorMapping.mapping = make(map[reflect.Type]map[reflect.Type]ConvertorFunc, 0)
	}

	toMap, foundToMap := convertorMapping.mapping[from]
	if !foundToMap {
		toMap = make(map[reflect.Type]ConvertorFunc, 0)
		convertorMapping.mapping[from] = toMap
	}

	toMap[to] = convertor

	return convertorMapping

}

func (convertorMapping *convertorMapping) getConvertor(from, to reflect.Type) (ConvertorFunc, bool) {

	if convertorMapping.mapping == nil {
		return nil, false
	}

	toMap, foundToMap := convertorMapping.mapping[from]
	if !foundToMap {
		return nil, false
	}

	convertor, convertorFound := toMap[to]
	return convertor, convertorFound

}

func (convertorMapping *convertorMapping) ConvertIn(in, out interface{}) error {
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		panic("Convert needs a pointer to convert to as output parameter")
	}

	if convertor, foundConvertor := convertorMapping.getConvertor(reflect.TypeOf(in), reflect.TypeOf(out).Elem()); foundConvertor {
		result, err := convertor(in)
		if err == nil {
			reflect.ValueOf(out).Elem().Set(reflect.ValueOf(result))
		}
		return err
	}

	return fmt.Errorf("no idea how to convert %v to %v", reflect.TypeOf(in), reflect.TypeOf(out).Elem())
}

func (convertorMapping *convertorMapping) ConvertTo(in interface{}, to reflect.Type) (interface{}, error) {
	if convertor, foundConvertor := convertorMapping.getConvertor(reflect.TypeOf(in), to); foundConvertor {
		result, err := convertor(in)
		return result, err
	}

	return nil, fmt.Errorf("no idea how to convert %v to %v", reflect.TypeOf(in), to)
}

// DefaultConvertor creates a connverter with sensitive defaults for the most basic data types
//
func DefaultConvertor() Convertor {
	convertor := &convertorMapping{}

	convertor.
		//
		// int to ... convertors
		//
		RegisterConvertor(reflect.TypeOf(0), reflect.TypeOf(""), func(in interface{}) (interface{}, error) {
			return strconv.Itoa(in.(int)), nil
		}).
		RegisterConvertor(reflect.TypeOf(0), reflect.TypeOf(0.0), func(in interface{}) (interface{}, error) {
			return float64(in.(int)), nil
		}).
		//
		// float to ... convertors
		//
		RegisterConvertor(reflect.TypeOf(0.0), reflect.TypeOf(""), func(in interface{}) (interface{}, error) {
			return fmt.Sprintf("%f", in), nil
		}).
		RegisterConvertor(reflect.TypeOf(0.0), reflect.TypeOf(0), func(in interface{}) (interface{}, error) {
			return int(math.Round(in.(float64))), nil
		}).
		//
		// string to ... convertors
		//
		RegisterConvertor(reflect.TypeOf(""), reflect.TypeOf(0), func(in interface{}) (interface{}, error) {
			i, err := strconv.ParseInt(in.(string), 10, 64)
			if err != nil {
				return 0, err
			}
			return int(i), err
		}).
		RegisterConvertor(reflect.TypeOf(""), reflect.TypeOf(0.0), func(in interface{}) (interface{}, error) {
			f, err := strconv.ParseFloat(in.(string), 64)
			if err != nil {
				return 0.0, err
			}
			return float64(f), err
		}).
		RegisterConvertor(reflect.TypeOf(""), reflect.TypeOf(time.Time{}), func(in interface{}) (interface{}, error) {
			t, err := iso8601.ParseString(in.(string))
			if err != nil {
				return time.Time{}, err
			}
			return t, err
		}).
		//
		// time.Time to ... convertors
		//
		RegisterConvertor(reflect.TypeOf(time.Time{}), reflect.TypeOf(""), func(in interface{}) (interface{}, error) {
			return in.(time.Time).Format(time.RFC3339), nil
		})

	return convertor
}
