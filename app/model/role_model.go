package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	FactionID   uint
	PhaseID     uint
	Name        string `gorm:"unique"`
	Priority    uint   `gorm:"check:priority > 0"`
	Score       int    `gorm:"default:1"`
	Quantity    int    `gorm:"default:1;check:score > 0"`
	Image       string `gorm:"type:text;default:''"`
	Description string `gorm:"type:text"`

	Faction Faction `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL"`
	Phase   Phase   `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL"`
}
