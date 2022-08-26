package state_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"uwwolf/module/game/state"
)

func TestIsShot(t *testing.T) {
	s := state.NewShotgun()

	//=============================================================
	// Before shooting
	assert.False(t, s.IsShot())

	//=============================================================
	// After shooting
	s.Shoot(1)

	assert.True(t, s.IsShot())
}

func TestShot(t *testing.T) {
	s := state.NewShotgun()

	//=============================================================
	// Target > 0
	assert.False(t, s.Shoot(0))

	//=============================================================
	// Shoot successfully
	assert.True(t, s.Shoot(1))

	//=============================================================
	// Shoot failed
	assert.False(t, s.Shoot(1))
	assert.False(t, s.Shoot(2))
}
