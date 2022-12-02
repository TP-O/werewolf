package model

import "uwwolf/game/enum"

type Player struct {
	ID       enum.PlayerID     `json:"id" gorm:"primaryKey;type:varchar(30)"`
	StatusID enum.PlayerStatus `json:"statusID" gorm:""`
	Status   PlayerStatus      `json:"status" gorm:"foreignKey:StatusId"`
}

func (Player) TableName() string {
	return "users"
}
