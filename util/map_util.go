package util

func ExistKeyInMap[K comparable, V any](m map[K]V, key K) bool {
	_, ok := m[key]

	return ok
}

func ExistValueInMap[K comparable, V comparable](m map[K]V, value V) bool {
	for _, v := range m {
		if v == value {
			return true
		}
	}

	return false
}
