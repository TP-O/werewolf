package model

import (
	"gorm.io/gorm"

	"uwwolf/app/types"
)

type Role struct {
	Id         types.RoleId        `gorm:"primaryKey" json:"id"`
	PhaseId    types.PhaseId       `gorm:"not null;uniqueIndex:idx_priority_in_phase" json:"phaseId"`
	Phase      Phase               `gorm:"foreignKey:PhaseId" json:"phase"`
	FactionId  types.FactionId     `gorm:"not null" json:"factionId"`
	Faction    Faction             `gorm:"foreignKey:FactionId" json:"faction"`
	IsDefault  bool                `gorm:"type:boolean;default:0" json:"isDefault"`
	Priority   int                 `gorm:"type:smallint;not null;uniqueIndex:idx_priority_in_phase;check:priority > 0" json:"priority"`
	Weight     int                 `gorm:"type:smallint;not null" json:"weight"`
	Set        int                 `gorm:"type:smallint;default:1;check:set = -1 AND set > 0" json:"set"`
	BeginRound types.RoundId       `gorm:"type:smallint;not null;check:begin_round > 0" json:"beginRound"`
	Expiration types.NumberOfTimes `gorm:"type:smallint;not null;check:expiration = -1 OR expiration > 0" json:"expiration"`
	DeletedAt  gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}
