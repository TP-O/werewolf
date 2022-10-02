package model

import "uwwolf/app/types"

type GameRecord struct {
	GameId   types.GameId   `gorm:"index;not null" json:"gameId"`
	Game     Game           `gorm:"foreignKey:GameId" json:"game"`
	PhaseId  types.PhaseId  `gorm:"not null" json:"phaseId"`
	Phase    Phase          `gorm:"foreignKey:PhaseId" json:"phase"`
	ActorId  types.PlayerId `gorm:"not null" json:"actorId"`
	Actor    Player         `gorm:"foreignKey:ActorId" json:"actor"`
	TargetId types.PlayerId `gorm:"not null" json:"targetId"`
	Target   Player         `gorm:"foreignKey:TargetId" json:"target"`
	ActionId types.ActionId `gorm:"not null" json:"actionId"`
	Action   Action         `gorm:"foreignKey:ActionId" json:"action"`
}
