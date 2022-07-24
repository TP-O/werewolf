package role

import (
	"time"
	"uwwolf/contract/itf"
	"uwwolf/contract/typ"
	"uwwolf/game/action"
)

func NewVillagerRole(game itf.IGame) *role {
	return &role{
		name: "Villager",
		game: game,
		skill: &itf.Skill{
			Action: action.NewVoteAction(2 * time.Second),
			Turn: &typ.SkillTurn{
				StartFrom:    2,
				NumberOfUses: -1,
			},
		},
	}
}
