package factory

import (
	"uwwolf/app/game/contract"
	"uwwolf/app/game/role"
	"uwwolf/app/types"
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

	default:
		return nil
	}
}
