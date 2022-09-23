package role

import (
	"uwwolf/app/game/action"
	"uwwolf/app/game/contract"
	"uwwolf/app/types"
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
			action:       action.NewVote(game, types.WerewolfFaction, player.Id(), 1),
			beginRoundId: setting.BeginRound,
			expiration:   setting.Expiration,
		},
	}
}
