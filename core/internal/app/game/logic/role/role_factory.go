package role

import (
	"fmt"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
)

func NewRole(id types.RoleID, world contract.World, playerID types.PlayerID) (contract.Role, error) {
	switch id {
	case constants.VillagerRoleID:
		return NewVillager(world, playerID)

	case constants.WerewolfRoleID:
		return NewWerewolf(world, playerID)

	case constants.HunterRoleID:
		return NewHunter(world, playerID)

	case constants.SeerRoleID:
		return NewSeer(world, playerID)

	case constants.TwoSistersRoleID:
		return NewTwoSister(world, playerID)

	default:
		return nil, fmt.Errorf("Non-existent role ¯\\_ಠ_ಠ_/¯")
	}
}
