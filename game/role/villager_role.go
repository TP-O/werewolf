package role

import "uwwolf/contract/itf"

func NewVillagerRole(game itf.IGame) *role {
	return &role{
		name: "Villager",
		game: game,
	}
}
