package role

import "uwwolf/contract/itf"

func NewAlphaWolfRole(game itf.IGame) *role {
	return &role{
		name: "Alpha Wolf",
		game: game,
	}
}
