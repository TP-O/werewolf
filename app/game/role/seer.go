package role

import (
	"uwwolf/app/game/action"
	"uwwolf/app/game/contract"
	"uwwolf/app/types"
)

func NewSeerRole(game contract.Game, setting *types.RoleSetting) contract.Role {
	return &role{
		id:        setting.Id,
		factionId: setting.FactionId,
		phaseId:   setting.PhaseId,
		game:      game,
		player:    game.Player(setting.OwnerId),
		skill: &skill{
			action:       action.NewProphecy(game),
			beginRoundId: setting.BeginRound,
			expiration:   setting.Expiration,
		},
	}
}
