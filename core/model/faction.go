package model

import "uwwolf/game/enum"

type Faction struct {
	ID   enum.FactionID `cql:"id" json:"id"`
	Name string         `cql:"name" json:"name"`
}
