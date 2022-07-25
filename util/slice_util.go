package util

import (
	"reflect"
)

func ExistElement[T any](arr []T, value T) bool {
	for _, el := range arr {
		if reflect.DeepEqual(el, value) {
			return true
		}
	}

	return false
}

func RemoveDuplicateElement[T any](collection []T) []T {
	newCollection := make([]T, 0)

	for _, e := range collection {
		if !ExistElement(newCollection, e) {
			newCollection = append(newCollection, e)
		}
	}

	return newCollection
}
