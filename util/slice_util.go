package util

import "github.com/samber/lo"

func Find[T comparable](collection []T, foundItem T) bool {
	_, ok := lo.Find(collection, func(item T) bool { return item == foundItem })

	return ok
}
