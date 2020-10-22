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
