package model

import "uwwolf/game/types"

type Action struct {
	ID   types.ActionID `cql:"id" json:"id"`
	Name string         `cql:"name" json:"name"`
}
