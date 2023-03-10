package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomInRange(t *testing.T) {
	min := 7
	max := 9
	for i := 0; i < 10; i++ {
		randNum := RandomInRange(min, max)

		assert.LessOrEqual(t, randNum, max)
		assert.GreaterOrEqual(t, randNum, min)
	}
}
func TestRandomIndex(t *testing.T) {
	//=============================================================
	// Test non-empty array
	nonEmptyArr := []int{1, 2, 3}
	for i := 0; i < 10; i++ {
		randIndex := RandomIndex(nonEmptyArr)

		assert.Less(t, randIndex, len(nonEmptyArr))
		assert.GreaterOrEqual(t, randIndex, 0)
	}

	//=============================================================
	// Test empty array
	emptyArr := []int{}
	randIndex := RandomIndex(emptyArr)

	assert.Equal(t, randIndex, -1)
}

func TestRandomElement(t *testing.T) {
	//=============================================================
	// Test non-empty array
	nonEmptyArr := []int{1, 2, 3}
	for i := 0; i < 10; i++ {
		randIndex, randElement := RandomElement(nonEmptyArr)

		assert.Less(t, randIndex, len(nonEmptyArr))
		assert.GreaterOrEqual(t, randIndex, 0)
		assert.Contains(t, nonEmptyArr, randElement)
	}

	//=============================================================
	// Test empty array
	emptyArr := []int{}
	randIndex, randElement := RandomElement(emptyArr)

	assert.Equal(t, randIndex, -1)
	assert.Equal(t, randElement, 0)
}
