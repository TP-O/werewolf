package factory

import (
	"uwwolf/game/contract"
	"uwwolf/game/enum"
	"uwwolf/game/role"
)

type roleFactory struct {
	roles map[uint]contract.Role
}

var roleFactoryInstance roleFactory

func init() {
	roleFactoryInstance = roleFactory{make(map[uint]contract.Role)}
}

func GetRoleFactory() *roleFactory {
	return &roleFactoryInstance
}

func (f *roleFactory) Create(key uint) contract.Role {
	if val, ok := f.roles[key]; ok {
		return val
	}

	switch key {
	case enum.VillagerRole:
		f.roles[key] = role.NewVillagerRole()

	case enum.WerewolfRole:
		f.roles[key] = role.NewWerewolfRole()

	case enum.HunterRole:
		f.roles[key] = role.NewHunterRole()

	case enum.SeerRole:
		f.roles[key] = role.NewSeerRole()

	case enum.TwoSistersRole:
		f.roles[key] = role.NewVillagerRole()

	case enum.AlphaWerewolfRole:
		f.roles[key] = role.NewWerewolfRole()

	default:
		return nil
	}

	return f.roles[key]
}
