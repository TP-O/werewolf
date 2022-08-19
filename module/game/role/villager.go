package role

import (
	"uwwolf/module/game/action"
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

const VillagerRoleName = "Villager"

func NewVillagerRole(game contract.Game, playerId types.PlayerId) contract.Role {
	return &role{
		id:      types.VillagerRole,
		phaseId: types.DayPhase,
		name:    VillagerRoleName,
		game:    game,
		player:  game.GetPlayer(playerId),
		skill: &skill{
			action:       action.NewVote(game, 1),
			numberOfUses: types.UnlimitedTimes,
			beginRoundId: types.FirstRound,
		},
	}
}
