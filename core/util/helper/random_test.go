package helper_test

import (
	"testing"
	"uwwolf/util/helper"

	"github.com/stretchr/testify/assert"
)

func TestRandomInRange(t *testing.T) {
	min := 7
	max := 9
	for i := 0; i < 10; i++ {
		randNum := helper.RandomInRange(min, max)

		assert.LessOrEqual(t, randNum, max)
		assert.GreaterOrEqual(t, randNum, min)
	}
}
func TestRandomIndex(t *testing.T) {
	//=============================================================
	// Test non-empty array
	nonEmptyArr := []int{1, 2, 3}
	for i := 0; i < 10; i++ {
		randIndex := helper.RandomIndex(nonEmptyArr)

		assert.Less(t, randIndex, len(nonEmptyArr))
		assert.GreaterOrEqual(t, randIndex, 0)
	}

	//=============================================================
	// Test empty array
	emptyArr := []int{}
	randIndex := helper.RandomIndex(emptyArr)

	assert.Equal(t, randIndex, -1)
}

func TestRandomElement(t *testing.T) {
	//=============================================================
	// Test non-empty array
	nonEmptyArr := []int{1, 2, 3}
	for i := 0; i < 10; i++ {
		randIndex, randElement := helper.RandomElement(nonEmptyArr)

		assert.Less(t, randIndex, len(nonEmptyArr))
		assert.GreaterOrEqual(t, randIndex, 0)
		assert.Contains(t, nonEmptyArr, randElement)
	}

	//=============================================================
	// Test empty array
	emptyArr := []int{}
	randIndex, randElement := helper.RandomElement(emptyArr)

	assert.Equal(t, randIndex, -1)
	assert.Equal(t, randElement, 0)
}
