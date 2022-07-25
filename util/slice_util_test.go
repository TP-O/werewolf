package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExistElement(t *testing.T) {
	// Comparable type
	arr1 := []int{1, 2, 3}

	assert.True(t, ExistElement(arr1, 2))
	assert.False(t, ExistElement(arr1, 5))

	// Uncomparable type
	type uncomparableType struct {
		val int
	}

	el1 := uncomparableType{val: 1}
	el2 := uncomparableType{val: 2}
	el3 := uncomparableType{val: 3}
	el4 := uncomparableType{val: 4}
	arr2 := []uncomparableType{el1, el2, el3}

	assert.True(t, ExistElement(arr2, el2))
	assert.False(t, ExistElement(arr2, el4))
}

func TestRemoveDuplicateElement(t *testing.T) {
	arr := []int{1, -9, 2, 7, 3, 9, 6, 4, 9, 1}
	newArr := RemoveDuplicateElement(arr)

	for i := 0; i < len(newArr)-1; i++ {
		for j := i + 1; j < len(newArr); j++ {
			assert.NotEqual(t, newArr[i], newArr[j])
		}
	}
}
