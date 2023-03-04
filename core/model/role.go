package model

import "uwwolf/game/types"

type Role struct {
	ID        types.RoleID    `cql:"id" json:"id"`
	FactionID types.FactionID `cql:"faction_id" json:"factionID"`
	PhaseID   types.PhaseID   `cql:"phase_id" json:"phaseID"`
	Name      string          `cql:"name" json:"name"`
}
