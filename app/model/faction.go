package model

import (
	"uwwolf/app/types"
)

type Faction struct {
	Id types.FactionId `gorm:"primaryKey;type:int" json:"id"`
}
