package model

import (
	"gorm.io/gorm"

	"uwwolf/app/types"
)

type Faction struct {
	Id        types.FactionId `gorm:"primaryKey;type:int" json:"id"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"deletedAt"`
}
