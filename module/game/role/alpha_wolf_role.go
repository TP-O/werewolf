package role

import (
	"uwwolf/module/game/core"
	"uwwolf/types"
)

const AlphaWolfRoleName = "AlphaWolf"

func NewAlphaWolfRole(game core.Game) *role {
	return &role{
		id:      types.AlphaWolfRole,
		phaseId: types.NightPhase,
		name:    AlphaWolfRoleName,
		game:    game,
	}
}
