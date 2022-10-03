package model

import (
	"uwwolf/app/types"
)

type Phase struct {
	Id types.PhaseId `gorm:"primaryKey;type:int" json:"id"`
}
