package role

import (
	"uwwolf/app/enum"
	"uwwolf/app/game/contract"
	"uwwolf/app/types"
)

func New(id types.RoleId, game contract.Game, playerId types.PlayerId) contract.Role {
	switch id {
	case enum.VillagerRoleId:
		return newVillager(game, playerId)

	case enum.WerewolfRoleId:
		return newWerewolf(game, playerId)

	case enum.HunterRoleId:
		return newHunter(game, playerId)

	case enum.SeerRoleId:
		return newSeer(game, playerId)

	case enum.TwoSistersRoleId:
		return newTwoSister(game, playerId)

	default:
		return nil
	}
}
