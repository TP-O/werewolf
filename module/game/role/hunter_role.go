package role

import (
	"uwwolf/contract/itf"
	"uwwolf/contract/typ"
	"uwwolf/game/action"
)

func NewHunterRole(game itf.IGame) *role {
	return &role{
		name: "Hunter",
		game: game,
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
