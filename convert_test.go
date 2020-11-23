package mappert

import (
	"reflect"
	"testing"
	"time"
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

func TestConvertToUsingIntToFloat(t *testing.T) {

	if result, err := DefaultConvertor().ConvertTo(42, reflect.TypeOf(0.0)); err != nil {
		t.Error(err)
	} else {
		if result != 42.0 {
			t.Error("expected to convert to 42.0, not", result)
		}
	}
}

func TestConvertToUsingStringToFloat(t *testing.T) {

	if result, err := DefaultConvertor().ConvertTo("42", reflect.TypeOf(0.0)); err != nil {
		t.Error(err)
	} else {
		if result != 42.0 {
			t.Error("expected to convert to 42.0, not", result)
		}
	}
}

func TestConvertToUsingFloatToInt(t *testing.T) {

	if result, err := DefaultConvertor().ConvertTo(42.49, reflect.TypeOf(0)); err != nil {
		t.Error(err)
	} else {
		if result != 42 {
			t.Error("expected to convert to 42, not", result)
		}
	}

	if result, err := DefaultConvertor().ConvertTo(42.50, reflect.TypeOf(0)); err != nil {
		t.Error(err)
	} else {
		if result != 43 {
			t.Error("expected to convert to 43, not", result)
		}
	}
}

func TestConvertToUsingStringToTimeAndViceVersa(t *testing.T) {

	if converted, err := DefaultConvertor().ConvertTo("2002-10-02T15:00:00Z", reflect.TypeOf(time.Time{})); err != nil {
		t.Error(err)
	} else {
		if result, err := DefaultConvertor().ConvertTo(converted, reflect.TypeOf("")); err != nil {
			t.Error(err)
		} else {
			if result != "2002-10-02T15:00:00Z" {
				t.Error("expected to convert to 2002-10-02T15:00:00Z, not", result)
			}
		}
	}
}
