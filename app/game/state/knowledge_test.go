package state_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"uwwolf/app/game/state"
	"uwwolf/app/types"
)

func TestIdentify(t *testing.T) {
	k := state.NewKnowledge()

	//=============================================================
	// Data is empty
	factionId := k.Identify("1")

	assert.Equal(t, factionId, types.UnknownFaction)

	//=============================================================
	// Data is provided
	k.Acquire("1", types.VillagerFaction)
	factionId = k.Identify("1")

	assert.Equal(t, factionId, types.VillagerFaction)
}

func TestAcquire(t *testing.T) {
	k := state.NewKnowledge()

	//=============================================================
	// New data
	res := k.Acquire("1", types.VillagerFaction)

	assert.True(t, res)

	//=============================================================
	// Old data
	playerId := types.PlayerId("2")
	k.Acquire(playerId, types.VillagerFaction)
	res = k.Acquire(playerId, types.VillagerFaction)

	assert.False(t, res)
}
