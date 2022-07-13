package util

import (
	"math/rand"
	"time"
)

func Intn(n int) int {
	seed()

	return rand.Intn(n)
}

func seed() {
	rand.Seed(time.Now().UnixNano())
}
