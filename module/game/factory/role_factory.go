package factory

import (
	"uwwolf/module/game"
	"uwwolf/module/game/role"
	"uwwolf/types"
)

type roleFactory struct {
	//
}

func (f *roleFactory) Create(key types.Role, game game.IGame) role.IRole {
	switch key {
	case types.VillagerRole:
		return role.NewVillagerRole(game)

	case types.WerewolfRole:
		return role.NewWerewolfRole(game)

	case types.HunterRole:
		return role.NewHunterRole(game)

	case types.SeerRole:
		return role.NewSeerRole(game)

	case types.TwoSistersRole:
		return role.NewVillagerRole(game)

	case types.AlphaWolfRole:
		return role.NewAlphaWolfRole(game)

	default:
		return nil
	}
}
