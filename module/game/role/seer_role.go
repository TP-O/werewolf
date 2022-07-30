package role

import (
	"uwwolf/module/game/action"
	"uwwolf/module/game/core"
	"uwwolf/types"
)

const SeerRoleName = "Seer"

func NewSeerRole(game core.Game) *role {
	return &role{
		id:      types.SeerRole,
		phaseId: types.NightPhase,
		name:    SeerRoleName,
		game:    game,
		activeSkill: &activeSkill{
			action:       action.NewProphecyAction(game),
			numberOfUses: types.UnlimitedUse,
			startRound:   types.FirstRound,
		},
	}
}
