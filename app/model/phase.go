package model

import (
	"gorm.io/gorm"

	"uwwolf/app/types"
)

type Phase struct {
	Id        types.PhaseId  `gorm:"primaryKey;type:int" json:"id"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}
