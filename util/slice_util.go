package util

import (
	"golang.org/x/exp/slices"
)

func ExistInArray[T comparable](arr []T, els ...T) bool {
	for _, el := range els {
		if !slices.Contains(arr, el) {
			return false
		}
	}

	return true
}
