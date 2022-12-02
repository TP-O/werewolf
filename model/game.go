package model

import (
	"database/sql"
	"uwwolf/game/enum"
)

type Game struct {
	ID               enum.GameID    `gorm:"primaryKey" json:"id"`
	WinningFactionID enum.FactionID `gorm:"" json:"factionID"`
	WinningFaction   Faction        `gorm:"foreignKey:WinningFactionId" json:"winingFaction"`
	StartedAt        sql.NullTime   `gorm:"index" json:"startedAt"`
	FinishedAt       sql.NullTime   `gorm:"index" json:"finishedAt"`
}
