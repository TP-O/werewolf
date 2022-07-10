package role

import (
	"uwwolf/game/action"
	"uwwolf/game/contract"
)

func NewHunterRole() *role {
	return &role{
		name: "Hunter",
		skill: &contract.Skill{
			Action: action.NewShootingAction(),
			Turn: &contract.SkillTurn{
				StartFrom:    2,
				NumberOfUses: -1,
			},
		},
		passive: &contract.Passive{
			Action: action.NewShootingAction(),
			Event: &contract.PassiveEvent{
				AfterDeath: true,
			},
		},
	}
}
