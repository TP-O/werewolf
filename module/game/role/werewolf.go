package role

import (
	"uwwolf/module/game/action"
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

const WerewolfRoleName = "Werewolf"

func NewWerewolfRole(game contract.Game, setting *types.RoleSetting) contract.Role {
	return &role{
		id:      types.WerewolfRole,
		phaseId: types.NightPhase,
		name:    WerewolfRoleName,
		game:    game,
		player:  game.GetPlayer(setting.OwnerId),
		skill: &skill{
			action:       action.NewVote(game, 1),
			beginRoundId: setting.BeginRound,
			expiration:   setting.Expiration,
		},
	}
}
