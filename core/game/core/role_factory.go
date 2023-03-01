package core

import (
	"fmt"
	"uwwolf/game/contract"
	"uwwolf/game/role"
	"uwwolf/game/types"
)

func NewRole(id types.RoleID, moderator contract.Game, playerID types.PlayerID) (contract.Role, error) {
	switch id {
	case role.VillagerRoleID:
		return role.NewVillager(moderator, playerID)

	case role.WerewolfRoleID:
		return role.NewWerewolf(moderator, playerID)

	case role.HunterRoleID:
		return role.NewHunter(moderator, playerID)

	case role.SeerRoleID:
		return role.NewSeer(moderator, playerID)

	case role.TwoSistersRoleID:
		return role.NewTwoSister(moderator, playerID)

	default:
		return nil, fmt.Errorf("Non-existent role ¯\\_ಠ_ಠ_/¯")
	}
}
