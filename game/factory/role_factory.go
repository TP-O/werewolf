package factory

import (
	"uwwolf/contract/itf"
	"uwwolf/enum"
	"uwwolf/game/role"
)

type roleFactory struct {
	//
}

var roleFactoryInstance itf.IFactory[uint, itf.IGame, itf.IRole]

func GetRoleFactory() itf.IFactory[uint, itf.IGame, itf.IRole] {
	return roleFactoryInstance
}

func (f *roleFactory) Create(key uint, game itf.IGame) itf.IRole {
	switch key {
	case enum.VillagerRole:
		return role.NewVillagerRole(game)

	case enum.WerewolfRole:
		return role.NewWerewolfRole(game)

	case enum.HunterRole:
		return role.NewHunterRole(game)

	case enum.SeerRole:
		return role.NewSeerRole(game)

	case enum.TwoSistersRole:
		return role.NewVillagerRole(game)

	case enum.AlphaWolfRole:
		return role.NewAlphaWolfRole(game)

	default:
		return nil
	}
}
