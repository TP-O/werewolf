package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExistKeyInMap(t *testing.T) {
	m := map[int]int{
		0: 1,
		1: 2,
		2: 3,
	}

	assert.True(t, ExistKeyInMap(m, 2))
	assert.False(t, ExistKeyInMap(m, 5))
}

func TestExistValueInMap(t *testing.T) {
	m := map[int]int{
		0: 1,
		1: 2,
		2: 3,
	}

	assert.True(t, ExistValueInMap(m, 2))
	assert.False(t, ExistValueInMap(m, 5))
}
