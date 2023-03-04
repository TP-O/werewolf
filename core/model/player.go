package model

import "uwwolf/game/types"

type Player struct {
	ID types.PlayerID `cql:"id" json:"id"`
}
