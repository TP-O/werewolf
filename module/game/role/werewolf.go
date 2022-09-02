package role

import (
	"uwwolf/module/game/action"
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

func NewWerewolfRole(game contract.Game, setting *types.RoleSetting) contract.Role {
	player := game.Player(setting.OwnerId)

	return &role{
		id:        setting.Id,
		factionId: setting.FactionId,
		phaseId:   setting.PhaseId,
		game:      game,
		player:    player,
		skill: &skill{
			action:       action.NewVote(game, player, 1),
			beginRoundId: setting.BeginRound,
			expiration:   setting.Expiration,
		},
	}
}
