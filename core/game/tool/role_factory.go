package game

import (
	"errors"
	"uwwolf/game/contract"
	rolee "uwwolf/game/role"
	"uwwolf/game/types"
)

func NewRole(id types.RoleID, game contract.Game, playerID types.PlayerID) (contract.Role, error) {
	switch id {
	case VillagerRoleID:
		return rolee.NewVillager(game, playerID)

	case WerewolfRoleID:
		return rolee.NewWerewolf(game, playerID)

	case HunterRoleID:
		return rolee.NewHunter(game, playerID)

	case SeerRoleID:
		return rolee.NewSeer(game, playerID)

	case TwoSistersRoleID:
		return rolee.NewTwoSister(game, playerID)

	default:
		return nil, errors.New("Non-existent role ¯\\_ಠ_ಠ_/¯")
	}
}
