package role

import (
	"uwwolf/game/action"
	"uwwolf/game/contract"
)

func NewSeerRole() *role {
	return &role{
		name: "Seer",
		skill: &contract.Skill{
			Action: action.NewProphecyAction(),
			Turn: &contract.SkillTurn{
				StartFrom:    2,
				NumberOfUses: -1,
			},
		},
	}
}
