package role

import (
	"uwwolf/module/game/action"
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

const VillagerRoleName = "Villager"

func NewVillagerRole(game contract.Game, setting *types.RoleSetting) contract.Role {
	return &role{
		id:      types.VillagerRole,
		phaseId: types.DayPhase,
		name:    VillagerRoleName,
		game:    game,
		player:  game.GetPlayer(setting.OwnerId),
		skill: &skill{
			action:       action.NewVote(game, 1),
			beginRoundId: setting.BeginRound,
			expiration:   setting.Expiration,
		},
	}
}
