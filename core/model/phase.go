package model

import "uwwolf/game/types"

type Phase struct {
	ID   types.PhaseID `cql:"id" json:"id"`
	Name string        `cql:"name" json:"name"`
}
