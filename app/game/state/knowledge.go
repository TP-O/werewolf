package state

import (
	"uwwolf/app/types"
	"uwwolf/app/util"
)

type Knowledge struct {
	// Map factions with their players
	Factions map[types.FactionId][]types.PlayerId `json:"factions"`

	// Map players with their faction
	Players map[types.PlayerId]types.FactionId `json:"players"`
}

func NewKnowledge() *Knowledge {
	return &Knowledge{
		Factions: make(map[types.FactionId][]types.PlayerId),
		Players:  make(map[types.PlayerId]types.FactionId),
	}
}

func (k *Knowledge) Identify(playerId types.PlayerId) types.FactionId {
	if util.ExistKeyInMap(k.Players, playerId) {
		return k.Players[playerId]
	}

	return types.UnknownFaction
}

func (k *Knowledge) Acquire(playerId types.PlayerId, factionId types.FactionId) bool {
	if k.Identify(playerId) != types.UnknownFaction {
		return false
	}

	k.Factions[factionId] = append(
		k.Factions[factionId],
		playerId,
	)
	k.Players[playerId] = factionId

	return true
}
