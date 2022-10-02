package model

import "uwwolf/app/types"

type RoleAction struct {
	RoleId     types.RoleId        `gorm:"primaryKey" json:"roleId"`
	Role       Role                `gorm:"foreignKey:RoleId" json:"role"`
	ActionId   types.ActionId      `gorm:"primaryKey" json:"actionId"`
	Action     Action              `gorm:"foreignKey:ActionId" json:"action"`
	Expiration types.NumberOfTimes `gorm:"type:smallint;not null;check:expiration = -1 OR expiration > 0" json:"expiration"`
}
