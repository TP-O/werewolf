package model

import "uwwolf/game/enum"

type Phase struct {
	ID enum.PhaseID `gorm:"primaryKey;type:int" json:"id"`
}
