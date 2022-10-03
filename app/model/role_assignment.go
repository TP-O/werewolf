package model

import (
	"uwwolf/app/types"
)

type RoleAssignment struct {
	GameId    types.GameId    `gorm:"index;uniqueIndex:idx_unique_assignment;not null" json:"gameId"`
	Game      Game            `gorm:"foreignKey:GameId" json:"game"`
	PlayerId  types.PlayerId  `gorm:"column:user_id;uniqueIndex:idx_unique_assignment;not null" json:"playerId"`
	Player    Player          `gorm:"foreignKey:PlayerId" json:"player"`
	RoleId    types.RoleId    `gorm:"uniqueIndex:idx_unique_assignment;not null" json:"roleId"`
	Role      Game            `gorm:"foreignKey:RoleId" json:"role"`
	FactionId types.FactionId `gorm:"" json:"factionId"`
	Faction   Faction         `gorm:"foreignKey:FactionId" json:"faction"`
}
