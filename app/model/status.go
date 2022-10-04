package model

import "uwwolf/app/types"

type Status struct {
	Id types.StatusId `gorm:"primaryKey;type:int" json:"id"`
}

func (Status) TableName() string {
	return "status"
}
