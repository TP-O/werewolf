package model

import (
	"database/sql"

	"uwwolf/app/types"
)

type Game struct {
	Id               types.GameId    `gorm:"primaryKey" json:"id"`
	WinningFactionId types.FactionId `gorm:"" json:"factionId"`
	WinningFaction   Faction         `gorm:"foreignKey:WinningFactionId" json:"winingFaction"`
	StartedAt        sql.NullTime    `gorm:"index" json:"startedAt"`
	FinishedAt       sql.NullTime    `gorm:"index" json:"finishedAt"`
}
