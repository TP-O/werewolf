package game

import (
	"fmt"
	"uwwolf/game/contract"
	"uwwolf/game/role"
	"uwwolf/game/types"
	"uwwolf/game/vars"
)

func NewRole(id types.RoleID, moderator contract.Game, playerID types.PlayerID) (contract.Role, error) {
	switch id {
	case vars.VillagerRoleID:
		return role.NewVillager(moderator, playerID)

	case vars.WerewolfRoleID:
		return role.NewWerewolf(moderator, playerID)

	case vars.HunterRoleID:
		return role.NewHunter(moderator, playerID)

	case vars.SeerRoleID:
		return role.NewSeer(moderator, playerID)

	case vars.TwoSistersRoleID:
		return role.NewTwoSister(moderator, playerID)

	default:
		return nil, fmt.Errorf("Non-existent role ¯\\_ಠ_ಠ_/¯")
	}
}
