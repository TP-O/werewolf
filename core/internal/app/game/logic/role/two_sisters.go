package role

import (
	"uwwolf/internal/app/game/logic/action"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
)

type twoSister struct {
	*role
}

func NewTwoSister(world contract.World, playerId types.PlayerId) (contract.Role, error) {
	return &twoSister{
		role: &role{
			id:           constants.TwoSistersRoleId,
			factionID:    constants.VillagerFactionId,
			phaseID:      constants.NightPhaseId,
			beginRoundID: constants.FirstRound,
			turnID:       constants.TwoSistersTurnID,
			world:        world,
			playerId:     playerId,
			abilities: []*ability{
				{
					action: action.NewRoleIdentify(
						world,
						constants.TwoSistersRoleId,
					),
					activeLimit: constants.Once,
				},
			},
		},
	}, nil
}
