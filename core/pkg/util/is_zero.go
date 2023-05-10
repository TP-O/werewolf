package util

import "reflect"

func IsZero(v any) bool {
	return reflect.ValueOf(v).IsZero()
}
