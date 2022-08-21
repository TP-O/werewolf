package role

import (
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

const AlphaWolfRoleName = "AlphaWolf"

type alphaWolf struct {
	role
}

func NewAlphaWolfRole(game contract.Game, setting *types.RoleSetting) contract.Role {
	return &role{
		id:      types.AlphaWolfRole,
		phaseId: types.NightPhase,
		name:    AlphaWolfRoleName,
		game:    game,
		player:  game.GetPlayer(setting.OwnerId),
		skill:   nil,
	}
}
