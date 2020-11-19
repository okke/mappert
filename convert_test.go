package mappert

import (
	"reflect"
	"testing"
)

func TestConvertInUsingIntToString(t *testing.T) {

	result := ""

	if err := DefaultConvertor().ConvertIn(42, &result); err != nil {
		t.Error(err)
	}

	if result != "42" {
		t.Error("expected to convert to 42, not", result)
	}
}

func TestConvertToUsingIntToString(t *testing.T) {

	if result, err := DefaultConvertor().ConvertTo(42, reflect.TypeOf("")); err != nil {
		t.Error(err)
	} else {
		if result != "42" {
			t.Error("expected to convert to 42, not", result)
		}
	}
}

func TestConvertInUsingStringToInt(t *testing.T) {

	result := 0

	if err := DefaultConvertor().ConvertIn("42", &result); err != nil {
		t.Error(err)
	}

	if result != 42 {
		t.Error("expected to convert to 42, not", result)
	}
}

func TestConvertInUsingInvalidStringToInt(t *testing.T) {

	result := 0

	if err := DefaultConvertor().ConvertIn("ab", &result); err == nil {
		t.Error("expected an error")
	}

}
