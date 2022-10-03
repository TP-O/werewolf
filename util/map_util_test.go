package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"uwwolf/util"
)

func TestExistKeyInMap(t *testing.T) {
	m := map[int]int{
		0: 1,
		1: 2,
		2: 3,
	}

	assert.True(t, util.ExistKeyInMap(m, 2))
	assert.False(t, util.ExistKeyInMap(m, 5))
}

func TestExistValueInMap(t *testing.T) {
	m := map[int]int{
		0: 1,
		1: 2,
		2: 3,
	}

	assert.True(t, util.ExistValueInMap(m, 2))
	assert.False(t, util.ExistValueInMap(m, 5))
}
