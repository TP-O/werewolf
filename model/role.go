package model

import "uwwolf/game/enum"

type Role struct {
	ID        enum.RoleID    `gorm:"primaryKey" json:"id"`
	PhaseID   enum.PhaseID   `gorm:"not null" json:"phaseID"`
	Phase     Phase          `gorm:"foreignKey:PhaseId" json:"phase"`
	FactionID enum.FactionID `gorm:"not null" json:"factionID"`
	Faction   Faction        `gorm:"foreignKey:FactionId" json:"faction"`
}
