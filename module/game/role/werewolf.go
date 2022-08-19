package role

import (
	"uwwolf/module/game/action"
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

const WerewolfRoleName = "Werewolf"

func NewWerewolfRole(game contract.Game, playerId types.PlayerId) contract.Role {
	return &role{
		id:      types.WerewolfRole,
		phaseId: types.NightPhase,
		name:    WerewolfRoleName,
		game:    game,
		player:  game.GetPlayer(playerId),
		skill: &skill{
			action:       action.NewVote(game, 1),
			numberOfUses: types.UnlimitedTimes,
			beginRoundId: types.FirstRound,
		},
	}
}
