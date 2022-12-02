package model

import "uwwolf/game/enum"

type Faction struct {
	ID enum.FactionID `gorm:"primaryKey;type:int" json:"id"`
}
