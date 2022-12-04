package model

import "uwwolf/game/enum"

type Role struct {
	ID        enum.RoleID    `cql:"id" json:"id"`
	FactionID enum.FactionID `cql:"faction_id" json:"factionID"`
	Name      string         `cql:"name" json:"name"`
}
