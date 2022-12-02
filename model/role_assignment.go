package model

import "uwwolf/game/enum"

type RoleAssignment struct {
	GameID    enum.GameID    `gorm:"index;uniqueIndex:idx_unique_assignment;not null" json:"gameID"`
	Game      Game           `gorm:"foreignKey:GameId" json:"game"`
	PlayerID  enum.PlayerID  `gorm:"column:user_id;uniqueIndex:idx_unique_assignment;not null" json:"playerID"`
	Player    Player         `gorm:"foreignKey:PlayerId" json:"player"`
	RoleID    enum.RoleID    `gorm:"uniqueIndex:idx_unique_assignment;not null" json:"roleID"`
	Role      Game           `gorm:"foreignKey:RoleId" json:"role"`
	FactionID enum.FactionID `gorm:"" json:"factionID"`
	Faction   Faction        `gorm:"foreignKey:FactionId" json:"faction"`
}
