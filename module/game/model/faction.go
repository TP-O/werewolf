package model

import (
	"uwwolf/types"

	"gorm.io/gorm"
)

type Faction struct {
	gorm.Model
	ID   types.PhaseId `gorm:"primarykey"`
	Name string        `gorm:"type:varchar(50);unique"`
}
