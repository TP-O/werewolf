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

func NewTwoSister(world contract.World, playerID types.PlayerID) (contract.Role, error) {
	return &twoSister{
		role: &role{
			id:           constants.TwoSistersRoleID,
			factionID:    constants.VillagerFactionID,
			phaseID:      constants.NightPhaseID,
			beginRoundID: constants.FirstRound,
			turnID:       constants.TwoSistersTurnID,
			world:        world,
			playerID:     playerID,
			abilities: []*ability{
				{
					action: action.NewRoleIdentify(
						world,
						constants.TwoSistersRoleID,
					),
					activeLimit: constants.Once,
				},
			},
		},
	}, nil
}
