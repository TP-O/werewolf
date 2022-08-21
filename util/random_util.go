package util

import (
	"math/rand"
	"time"
)

func RandomInRange(max int, min int) int {
	seed()

	return rand.Intn(max-min+1) + min
}

func RandomIndex[T any](arr []T) int {
	seed()

	if len(arr) == 0 {
		return -1
	}

	return rand.Intn(len(arr))
}

func RandomElement[T any](arr []*T) (*T, int) {
	seed()

	randomIndex := RandomIndex(arr)

	if randomIndex == -1 {
		return nil, -1
	}

	return arr[randomIndex], randomIndex
}

func seed() {
	rand.Seed(time.Now().UnixNano())
}
