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

	return rand.Intn(len(arr))
}

func seed() {
	rand.Seed(time.Now().UnixNano())
}
