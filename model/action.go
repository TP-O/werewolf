package model

import "uwwolf/game/enum"

type Action struct {
	ID   enum.ActionID `cql:"id" json:"id"`
	Name string        `cql:"name" json:"name"`
}
