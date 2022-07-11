package role

import (
	"uwwolf/contract/itf"
	"uwwolf/contract/typ"
	"uwwolf/game/action"
)

func NewHunterRole() *role {
	return &role{
		name: "Hunter",
		skill: &itf.Skill{
			Action: action.NewShootingAction(),
			Turn: &typ.SkillTurn{
				StartFrom:    2,
				NumberOfUses: -1,
			},
		},
		passive: &itf.Passive{
			Action: action.NewShootingAction(),
			Event: &typ.PassiveEvent{
				AfterDeath: true,
			},
		},
	}
}
