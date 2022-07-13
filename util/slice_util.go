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
