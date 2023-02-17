package model

import "uwwolf/game/enum"

type PlayerStatus struct {
	ID   enum.PlayerStatus `cql:"id" json:"id"`
	Name string            `cql:"name" json:"name"`
}
