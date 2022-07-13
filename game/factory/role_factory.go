package factory

import (
	"uwwolf/contract/itf"
	"uwwolf/enum"
	"uwwolf/game/role"
)

type roleFactory struct {
	//
}

var roleFactoryInstance itf.IFactory[uint, itf.IRole]

func GetRoleFactory() itf.IFactory[uint, itf.IRole] {
	return roleFactoryInstance
}

func (f *roleFactory) Create(key uint) itf.IRole {
	switch key {
	case enum.VillagerRole:
		return role.NewVillagerRole()

	case enum.WerewolfRole:
		return role.NewWerewolfRole()

	case enum.HunterRole:
		return role.NewHunterRole()

	case enum.SeerRole:
		return role.NewSeerRole()

	case enum.TwoSistersRole:
		return role.NewVillagerRole()

	case enum.AlphaWolfRole:
		return role.NewAlphaWolfRole()

	default:
		return nil
	}
}
