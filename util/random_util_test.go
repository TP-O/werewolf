package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomInRange(t *testing.T) {
	max := 9
	min := 7

	for i := 0; i < 10; i++ {
		randNum := RandomInRange(max, min)

		assert.LessOrEqual(t, randNum, max)
		assert.GreaterOrEqual(t, randNum, min)
	}
}
func TestRandomIndex(t *testing.T) {
	arr := []int{1, 2, 3}

	for i := 0; i < 10; i++ {
		randIndex := RandomIndex(arr)

		assert.Less(t, randIndex, len(arr))
		assert.GreaterOrEqual(t, randIndex, 0)
	}
}
