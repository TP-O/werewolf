package model

import (
	"uwwolf/app/types"
)

type Player struct {
	Id       types.PlayerId `json:"id" gorm:"primaryKey;type:varchar(30)"`
	StatusId types.StatusId `json:"statusId" gorm:""`
	Status   Status         `json:"status" gorm:"foreignKey:StatusId"`
}

func (Player) TableName() string {
	return "users"
}
