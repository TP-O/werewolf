package model

import "uwwolf/game/types"

type Faction struct {
	ID   types.FactionID `cql:"id" json:"id"`
	Name string          `cql:"name" json:"name"`
}
