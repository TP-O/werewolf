package role

import (
	"fmt"
	"uwwolf/internal/app/game/logic/declare"
	"uwwolf/internal/app/game/logic/mechanism/contract"
	"uwwolf/internal/app/game/logic/types"
)

func NewRole(id types.RoleID, world contract.World, playerID types.PlayerID) (contract.Role, error) {
	switch id {
	case declare.VillagerRoleID:
		return NewVillager(world, playerID)

	case declare.WerewolfRoleID:
		return NewWerewolf(world, playerID)

	case declare.HunterRoleID:
		return NewHunter(world, playerID)

	case declare.SeerRoleID:
		return NewSeer(world, playerID)

	case declare.TwoSistersRoleID:
		return NewTwoSister(world, playerID)

	default:
		return nil, fmt.Errorf("Non-existent role ¯\\_ಠ_ಠ_/¯")
	}
}
