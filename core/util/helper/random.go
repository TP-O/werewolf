package helper

import (
	"math/rand"
	"time"
)

// RandomInRange returns a random integer number from min to max.
func RandomInRange(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

// RandomIndex returns a random index of the given array.
// Returns -1 if the array is empty.
func RandomIndex[T any](arr []T) int {
	if len(arr) == 0 {
		return -1
	}

	return RandomInRange(0, len(arr)-1)
}

// RandomElement returns a random index and its value of the given array.
// Returns -1 and zero value if the array is empty.
func RandomElement[T any](arr []T) (int, T) {
	if randomIndex := RandomIndex(arr); randomIndex == -1 {
		var zero T
		return -1, zero
	} else {
		return randomIndex, arr[randomIndex]
	}
}
