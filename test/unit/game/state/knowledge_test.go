package state_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"uwwolf/module/game/state"
	"uwwolf/types"
)

func TestIdentify(t *testing.T) {
	k := state.NewKnowledge()

	//=============================================================
	// Data is empty
	factionId := k.Identify(1)

	assert.Equal(t, factionId, types.UnknownFaction)

	//=============================================================
	// Data is provided
	k.Acquire(1, types.VillagerFaction)
	factionId = k.Identify(1)

	assert.Equal(t, factionId, types.VillagerFaction)
}

func Acquire(t *testing.T) {
	k := state.NewKnowledge()

	//=============================================================
	// New data
	assert.True(t, k.Acquire(1, types.VillagerFaction))

	//=============================================================
	// Old data
	k.Acquire(2, types.VillagerFaction)

	assert.False(t, k.Acquire(2, types.VillagerFaction))
}
