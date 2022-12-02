package model

import "uwwolf/game/enum"

type GameRecord struct {
	GameID   enum.GameID   `gorm:"index;not null" json:"gameID"`
	Game     Game          `gorm:"foreignKey:GameId" json:"game"`
	PhaseID  enum.PhaseID  `gorm:"not null" json:"phaseID"`
	Phase    Phase         `gorm:"foreignKey:PhaseId" json:"phase"`
	ActorID  enum.PlayerID `gorm:"not null" json:"actorID"`
	Actor    Player        `gorm:"foreignKey:ActorId" json:"actor"`
	TargetID enum.PlayerID `gorm:"not null" json:"targetID"`
	Target   Player        `gorm:"foreignKey:TargetId" json:"target"`
	ActionID enum.ActionID `gorm:"not null" json:"actionID"`
	Action   Action        `gorm:"foreignKey:ActionId" json:"action"`
}
