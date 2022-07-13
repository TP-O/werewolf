package model

import "gorm.io/gorm"

type Faction struct {
	gorm.Model
	Name string `gorm:"unique"`

	Roles []Role
}
