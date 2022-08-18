package state

import (
	"uwwolf/types"
	"uwwolf/util"
)

type Prediction struct {
	// Map factions with their players
	factions map[types.FactionId][]types.PlayerId

	// Map players with their faction
	players map[types.PlayerId]types.FactionId
}

func NewPrediction() *Prediction {
	return &Prediction{
		factions: make(map[types.FactionId][]types.PlayerId),
		players:  make(map[types.PlayerId]types.FactionId),
	}
}

func (p *Prediction) Identify(playerId types.PlayerId) types.FactionId {
	if util.ExistKeyInMap(p.players, playerId) {
		return p.players[playerId]
	}

	return types.UnknownFaction
}

func (p *Prediction) Add(faction types.FactionId, player types.PlayerId) {
	p.factions[faction] = append(
		p.factions[faction],
		player,
	)
	p.players[player] = faction
}
