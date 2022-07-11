package role

import (
	"uwwolf/contract/itf"
	"uwwolf/contract/typ"
	"uwwolf/game/action"
)

func NewSeerRole() *role {
	return &role{
		name: "Seer",
		skill: &itf.Skill{
			Action: action.NewProphecyAction(),
			Turn: &typ.SkillTurn{
				StartFrom:    2,
				NumberOfUses: -1,
			},
		},
	}
}
