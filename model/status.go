package model

import "uwwolf/game/enum"

type PlayerStatus struct {
	ID enum.PlayerStatus `gorm:"primaryKey;type:int" json:"id"`
}

func (PlayerStatus) TableName() string {
	return "player_status"
}
