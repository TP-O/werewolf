package model

import "gorm.io/gorm"

type Phase struct {
	gorm.Model
	Name string `gorm:"unique"`

	Roles []Role
}
