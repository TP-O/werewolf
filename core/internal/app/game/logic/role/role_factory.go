package role

import (
	"fmt"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
)

func NewRole(id types.RoleId, world contract.World, playerID types.PlayerId) (contract.Role, error) {
	switch id {
	case constants.VillagerRoleId:
		return NewVillager(world, playerID)

	case constants.WerewolfRoleId:
		return NewWerewolf(world, playerID)

	case constants.HunterRoleId:
		return NewHunter(world, playerID)

	case constants.SeerRoleId:
		return NewSeer(world, playerID)

	case constants.TwoSistersRoleId:
		return NewTwoSister(world, playerID)

	default:
		return nil, fmt.Errorf("Non-existent role ¯\\_ಠ_ಠ_/¯")
	}
}
