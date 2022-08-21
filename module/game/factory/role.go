package factory

import (
	"uwwolf/module/game/contract"
	"uwwolf/module/game/role"
	"uwwolf/types"
)

type roleFactory struct {
	//
}

func (f *roleFactory) Create(id types.RoleId, game contract.Game, setting *types.RoleSetting) contract.Role {
	switch id {
	case types.VillagerRole:
		return role.NewVillagerRole(game, setting)

	case types.WerewolfRole:
		return role.NewWerewolfRole(game, setting)

	case types.HunterRole:
		return role.NewHunterRole(game, setting)

	case types.SeerRole:
		return role.NewSeerRole(game, setting)

	case types.TwoSistersRole:
		return role.NewTwoSisterRole(game, setting)

	case types.AlphaWolfRole:
		return role.NewAlphaWolfRole(game, setting)

	default:
		return nil
	}
}
