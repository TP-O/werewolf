package role

import (
	"uwwolf/module/game/action"
	"uwwolf/module/game/core"
	"uwwolf/types"
)

const HunterRoleName = "Hunter"

func NewHunterRole(game core.Game) *role {
	return &role{
		id:      types.HunterRole,
		phaseId: types.DayPhase,
		name:    HunterRoleName,
		game:    game,
		passiveSkill: &passiveSkill{
			action:     action.NewShootingAction(game),
			afterDeath: true,
		},
	}
}
