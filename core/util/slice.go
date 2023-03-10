package util

func FilterSlice[T any](slice []T, predicate func(T) bool) []T {
	result := []T{}
	for _, item := range slice {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func IsDuplicateSlice[T comparable](slice []T) bool {
	seen := make(map[T]bool)
	for _, item := range slice {
		if seen[item] {
			return true
		}
		seen[item] = true
	}
	return false
}
