package model

import (
	"uwwolf/types"

	"gorm.io/gorm"
)

type Phase struct {
	gorm.Model
	ID   types.RoleId `gorm:"primarykey"`
	Name string       `gorm:"type:varchar(50);unique"`
}
