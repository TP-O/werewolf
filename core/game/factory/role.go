package factory

import (
	"errors"
	"uwwolf/game/contract"
	"uwwolf/game/core/role"
	"uwwolf/game/enum"
)

func NewRole(id enum.RoleID, game contract.Game, playerID enum.PlayerID) (contract.Role, error) {
	switch id {
	case enum.VillagerRoleID:
		return role.NewVillager(game, playerID)

	case enum.WerewolfRoleID:
		return role.NewWerewolf(game, playerID)

	case enum.HunterRoleID:
		return role.NewHunter(game, playerID)

	case enum.SeerRoleID:
		return role.NewSeer(game, playerID)

	case enum.TwoSistersRoleID:
		return role.NewTwoSister(game, playerID)

	default:
		return nil, errors.New("Non-existent role ¯\\_ಠ_ಠ_/¯")
	}
}
