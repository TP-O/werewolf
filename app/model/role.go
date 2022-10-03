package model

import (
	"uwwolf/app/types"
)

type Role struct {
	Id        types.RoleId    `gorm:"primaryKey" json:"id"`
	PhaseId   types.PhaseId   `gorm:"not null;uniqueIndex:idx_priority_in_phase" json:"phaseId"`
	Phase     Phase           `gorm:"foreignKey:PhaseId" json:"phase"`
	FactionId types.FactionId `gorm:"not null" json:"factionId"`
	Faction   Faction         `gorm:"foreignKey:FactionId" json:"faction"`
}
