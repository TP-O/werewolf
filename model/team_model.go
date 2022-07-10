package model

import "gorm.io/gorm"

type Team struct {
	gorm.Model
	Name string `gorm:"unique"`

	Roles []Role
}
