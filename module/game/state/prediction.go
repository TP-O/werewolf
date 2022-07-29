package state

import (
	"uwwolf/types"
	"uwwolf/util"
)

type Prediction struct {
	factions map[types.Faction][]int
	players  map[int]types.Faction
}

func NewPrediction() Prediction {
	return Prediction{
		factions: make(map[types.Faction][]int),
		players:  make(map[int]types.Faction),
	}
}

func (p *Prediction) Identify(player int) types.Faction {
	if util.ExistKeyInMap(p.players, player) {
		return p.players[player]
	}

	return types.UnknownFaction
}

func (p *Prediction) Add(faction types.Faction, player int) {
	p.factions[faction] = append(
		p.factions[faction],
		player,
	)
	p.players[player] = faction
}

func (p *Prediction) Knowledge() map[types.Faction][]int {
	return p.factions
}
