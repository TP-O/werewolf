package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"uwwolf/util"
)

func TestRandomInRange(t *testing.T) {
	max := 9
	min := 7

	for i := 0; i < 10; i++ {
		randNum := util.RandomInRange(max, min)

		assert.LessOrEqual(t, randNum, max)
		assert.GreaterOrEqual(t, randNum, min)
	}
}
func TestRandomIndex(t *testing.T) {
	//=============================================================
	// Test non-empty array
	nonEmptyArr := []int{1, 2, 3}

	for i := 0; i < 10; i++ {
		randIndex := util.RandomIndex(nonEmptyArr)

		assert.Less(t, randIndex, len(nonEmptyArr))
		assert.GreaterOrEqual(t, randIndex, 0)
	}

	//=============================================================
	// Test empty array
	emptyArr := []int{}

	randIndex := util.RandomIndex(emptyArr)

	assert.Equal(t, randIndex, -1)
}

func TestRandomElement(t *testing.T) {
	//=============================================================
	// Test non-empty array
	nonEmptyArr := []int{1, 2, 3}

	for i := 0; i < 10; i++ {
		randIndex, randElement := util.RandomElement(nonEmptyArr)

		assert.Less(t, randIndex, len(nonEmptyArr))
		assert.GreaterOrEqual(t, randIndex, 0)
		assert.Contains(t, nonEmptyArr, randElement)
	}

	//=============================================================
	// Test empty array
	emptyArr := []int{}

	randIndex, randElement := util.RandomElement(emptyArr)

	assert.Equal(t, randIndex, -1)
	assert.Equal(t, randElement, 0)
}
