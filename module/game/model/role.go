// package model

// import (
// 	"uwwolf/types"

// 	"gorm.io/gorm"
// )

// type Role struct {
// 	gorm.Model
// 	FactionID   types.FactionId
// 	PhaseID     types.PhaseId
// 	Name        string `gorm:"unique"`
// 	Priority    int    `gorm:"check:priority > 0"`
// 	Score       int    `gorm:"default:1"`
// 	Quantity    int    `gorm:"default:1;check:score > 0"`
// 	Image       string `gorm:"type:text;default:''"`
// 	Description string `gorm:"type:text"`

// 	Faction Faction `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL"`
// 	Phase   Phase   `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL"`
// }

package model

import (
	"uwwolf/types"

	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	ID        types.RoleId  `gorm:"primarykey"`
	PhaseID   types.PhaseId `gorm:"uniqueIndex:idx_priority_in_phase"`
	FactionID types.FactionId
	Name      string `gorm:"type:varchar(50);unique"`
	Priority  int    `gorm:"type:integer;uniqueIndex:idx_priority_in_phase;check:priority > 0"`
	Weight    int    `gorm:"type:integer;default:0"`
	Set       int    `gorm:"type:integer;default:1;check:set > -2 and set <> 0"`

	Phase   Phase   `gorm:""`
	Faction Faction `gorm:""`
}
