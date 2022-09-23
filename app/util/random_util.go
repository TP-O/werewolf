package util

import (
	"math/rand"
	"time"
)

func RandomInRange(max int, min int) int {
	randomSeed()

	return rand.Intn(max-min+1) + min
}

func RandomIndex[T any](arr []T) int {
	randomSeed()

	if len(arr) == 0 {
		return -1
	}

	return rand.Intn(len(arr))
}

func RandomElement[T any](arr []T) (int, T) {
	randomSeed()

	randomIndex := RandomIndex(arr)

	if randomIndex == -1 {
		var zero T

		return -1, zero
	}

	return randomIndex, arr[randomIndex]
}

func randomSeed() {
	rand.Seed(time.Now().UnixNano())
}
