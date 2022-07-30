package role

import (
	"uwwolf/module/game/action"
	"uwwolf/module/game/core"
	"uwwolf/types"
)

const WerewolfRoleName = "Werewolf"

func NewWerewolfRole(game core.Game) *role {
	return &role{
		id:      types.WerewolfRole,
		phaseId: types.NightPhase,
		name:    WerewolfRoleName,
		game:    game,
		activeSkill: &activeSkill{
			action:       action.NewVoteAction(game),
			numberOfUses: types.UnlimitedUse,
			startRound:   types.FirstRound,
		},
	}
}
