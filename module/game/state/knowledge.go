package state

import (
	"uwwolf/types"
	"uwwolf/util"
)

type Knowledge struct {
	// Map factions with their players
	factions map[types.FactionId][]types.PlayerId

	// Map players with their faction
	players map[types.PlayerId]types.FactionId
}

func NewKnowledge() *Knowledge {
	return &Knowledge{
		factions: make(map[types.FactionId][]types.PlayerId),
		players:  make(map[types.PlayerId]types.FactionId),
	}
}

func (k *Knowledge) Identify(playerId types.PlayerId) types.FactionId {
	if util.ExistKeyInMap(k.players, playerId) {
		return k.players[playerId]
	}

	return types.UnknownFaction
}

func (k *Knowledge) Acquire(playerId types.PlayerId, factionId types.FactionId) bool {
	if k.Identify(playerId) != types.UnknownFaction {
		return false
	}

	k.factions[factionId] = append(
		k.factions[factionId],
		playerId,
	)
	k.players[playerId] = factionId

	return true
}
