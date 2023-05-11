package role

import (
	"fmt"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
)

func NewRole(id types.RoleId, moderator contract.Moderator, playerID types.PlayerId) (contract.Role, error) {
	switch id {
	case constants.VillagerRoleId:
		return NewVillager(moderator, playerID)

	case constants.WerewolfRoleId:
		return NewWerewolf(moderator, playerID)

	case constants.HunterRoleId:
		return NewHunter(moderator, playerID)

	case constants.SeerRoleId:
		return NewSeer(moderator, playerID)

	case constants.TwoSistersRoleId:
		return NewTwoSister(moderator, playerID)

	default:
		return nil, fmt.Errorf("Non-existent role ¯\\_ಠ_ಠ_/¯")
	}
}
