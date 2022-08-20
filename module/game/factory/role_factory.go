package factory

import (
	"uwwolf/module/game/contract"
	"uwwolf/module/game/role"
	"uwwolf/types"
)

type roleFactory struct {
	//
}

func (f *roleFactory) Create(key types.RoleId, game contract.Game, playerId types.PlayerId) contract.Role {
	switch key {
	case types.VillagerRole:
		return role.NewVillagerRole(game, playerId)

	case types.WerewolfRole:
		return role.NewWerewolfRole(game, playerId)

	case types.HunterRole:
		return role.NewHunterRole(game, playerId)

	case types.SeerRole:
		return role.NewSeerRole(game, playerId)

	case types.TwoSistersRole:
		return role.NewVillagerRole(game, playerId)

	case types.AlphaWolfRole:
		return role.NewAlphaWolfRole(game, playerId)

	default:
		return nil
	}
}
