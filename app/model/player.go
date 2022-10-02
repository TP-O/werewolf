package model

import (
	"uwwolf/app/types"
)

type Player struct {
	Id types.PlayerId `gorm:"primaryKey;type:varchar(30)" json:"id"`
}

func (Player) TableName() string {
	return "users"
}
