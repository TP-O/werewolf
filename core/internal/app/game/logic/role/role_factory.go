package role

import (
	"errors"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
)

// RoleFactory creates role instance.
type RoleFactory struct {
	//
}

var _ contract.RoleFactory = (*RoleFactory)(nil)

// CreateById creates a role with the given ID.
func (rf RoleFactory) CreateById(
	id types.RoleId,
	moderator contract.Moderator,
	playerID types.PlayerId,
) (contract.Role, error) {
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
		return nil, errors.New("Non-existent role ¯\\_ಠ_ಠ_/¯")
	}
}
