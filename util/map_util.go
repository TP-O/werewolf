package util

import "github.com/samber/lo"

func FindKey[K comparable, V any](mapping map[K]V, foundKey K) bool {
	return Find(lo.Keys(mapping), foundKey)
}

func FindValue[K comparable, V comparable](mapping map[K]V, foundValue V) bool {
	return Find(lo.Values(mapping), foundValue)
}
