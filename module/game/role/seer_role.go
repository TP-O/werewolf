package role

import (
	"uwwolf/contract/itf"
	"uwwolf/contract/typ"
	"uwwolf/game/action"
)

func NewSeerRole(game itf.IGame) *role {
	return &role{
		name: "Seer",
		game: game,
		skill: &itf.Skill{
			Action: action.NewProphecyAction(),
			Turn: &typ.SkillTurn{
				StartFrom:    2,
				NumberOfUses: -1,
			},
		},
	}
}
