package role

import (
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

const AlphaWolfRoleName = "Alpha Wolf"

func NewAlphaWolfRole(game contract.Game, setting *types.RoleSetting) contract.Role {
	return &role{
		id:        types.AlphaWolfRole,
		factionId: setting.FactionId,
		phaseId:   types.NightPhase,
		name:      AlphaWolfRoleName,
		game:      game,
		player:    game.Player(setting.OwnerId),
		skill: &skill{
			action:       nil,
			beginRoundId: setting.BeginRound,
			expiration:   setting.Expiration,
		},
	}
}
