package model

import "uwwolf/game/enum"

type Action struct {
	ID enum.ActionID `gorm:"primaryKey;type:int" json:"id"`
}
