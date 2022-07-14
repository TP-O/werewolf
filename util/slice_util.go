package util

import (
	"reflect"

	"github.com/samber/lo"
)

func Find[T comparable](collection []T, foundItem T) bool {
	_, ok := lo.Find(collection, func(item T) bool { return item == foundItem })

	return ok
}

func DeepFind[T any](collection []T, foundItem T) bool {
	_, ok := lo.Find(collection, func(item T) bool { return reflect.DeepEqual(item, foundItem) })

	return ok
}

func RemoveDuplicate[T any](collection []T) []T {
	newCollection := make([]T, 0)

	for _, e := range collection {
		if !DeepFind(newCollection, e) {
			newCollection = append(newCollection, e)
		}
	}

	return newCollection
}
