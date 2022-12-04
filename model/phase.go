package model

import "uwwolf/game/enum"

type Phase struct {
	ID   enum.PhaseID `cql:"id" json:"id"`
	Name string       `cql:"name" json:"name"`
}
