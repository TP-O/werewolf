package role

import (
	"uwwolf/module/game/action"
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

func NewTwoSisterRole(game contract.Game, setting *types.RoleSetting) contract.Role {
	return &role{
		id:        setting.Id,
		factionId: setting.FactionId,
		phaseId:   setting.PhaseId,
		game:      game,
		player:    game.Player(setting.OwnerId),
		skill: &skill{
			action:       action.NewRecognition(game, setting.Id),
			beginRoundId: setting.BeginRound,
			expiration:   setting.Expiration,
		},
	}
}
