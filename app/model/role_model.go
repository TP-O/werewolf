package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	TeamID      uint
	Name        string `gorm:"unique"`
	Score       int    `gorm:"default:1;check:score > 0"`
	Quantity    int    `gorm:"default:1;check:score > 0"`
	Image       string `gorm:"type:text;default:''"`
	Description string `gorm:"type:text"`

	Team Team `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL"`
}
