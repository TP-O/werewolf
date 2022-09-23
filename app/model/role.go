package model

import (
	"time"

	"gorm.io/gorm"

	"uwwolf/app/types"
)

type Role struct {
	Id         types.RoleId        `gorm:"type:uuid;default:uuid_generate_v4()" json:"id"`
	PhaseId    types.PhaseId       `gorm:"uniqueIndex:idx_priority_in_phase" json:"phaseId"`
	FactionId  types.FactionId     `gorm:"" json:"factionId"`
	IsDefault  bool                `gorm:"type:boolean;default=false" json:"isDefault"`
	Priority   int                 `gorm:"type:smallint;uniqueIndex:idx_priority_in_phase;check:priority >= 0" json:"priority"`
	Weight     int                 `gorm:"type:smallint;default:1"  json:"weight"`
	Set        int                 `gorm:"type:smallint;default:1;check:set >= -1 and set <> 0" json:"set"`
	BeginRound types.RoundId       `gorm:"type:smallint;default:1;check:begin_round >= 1" json:"beginRound"`
	Expiration types.NumberOfTimes `gorm:"type:smallint;default:1;check:set >= -1 and set <> 0" json:"expiration"`
	CreatedAt  time.Time           `gorm:"" json:"createdAt"`
	UpdatedAt  time.Time           `gorm:"" json:"updatedAt"`
	DeletedAt  gorm.DeletedAt      `gorm:"index" json:"deletedAt"`
}
