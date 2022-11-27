package factory

import (
	"uwwolf/app/game/config"
	"uwwolf/app/game/contract"
	"uwwolf/app/game/core/role"
	"uwwolf/app/game/types"
)

func NewRole(id types.RoleID, game contract.Game, playerID types.PlayerID) contract.Role {
	switch id {
	case config.VillagerRoleID:
		return role.NewVillager(game, playerID)

	case config.WerewolfRoleID:
		return role.NewWerewolf(game, playerID)

	case config.HunterRoleID:
		return role.NewHunter(game, playerID)

	case config.SeerRoleID:
		return role.NewSeer(game, playerID)

	case config.TwoSistersRoleID:
		return role.NewTwoSister(game, playerID)

	default:
		return nil
	}
}
