package role

import (
	"uwwolf/module/game/action"
	"uwwolf/module/game/core"
	"uwwolf/types"
)

const VillagerRoleName = "Villager"

func NewVillagerRole(game core.Game) Role {
	return &role{
		id:      types.VillagerRole,
		phaseId: types.DayPhase,
		name:    VillagerRoleName,
		game:    game,
		activeSkill: &activeSkill{
			action:       action.NewVoteAction(game),
			numberOfUses: types.UnlimitedUse,
			startRound:   types.FirstRound,
		},
	}
}
