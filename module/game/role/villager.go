package role

import (
	"uwwolf/module/game/action"
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

const VillagerRoleName = "Villager"

func NewVillagerRole(game contract.Game, setting *types.RoleSetting) contract.Role {
	player := game.Player(setting.OwnerId)

	return &role{
		id:        types.VillagerRole,
		factionId: setting.FactionId,
		phaseId:   types.DayPhase,
		name:      VillagerRoleName,
		game:      game,
		player:    player,
		skill: &skill{
			action:       action.NewVote(game, player, 1),
			beginRoundId: setting.BeginRound,
			expiration:   setting.Expiration,
		},
	}
}
