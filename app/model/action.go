package model

import (
	"uwwolf/app/types"
)

type Action struct {
	Id types.ActionId `gorm:"primaryKey;type:int" json:"id"`
}
