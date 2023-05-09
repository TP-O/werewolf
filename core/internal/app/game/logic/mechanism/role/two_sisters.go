package role

import (
	"uwwolf/internal/app/game/logic/declare"
	"uwwolf/internal/app/game/logic/mechanism/action"
	"uwwolf/internal/app/game/logic/mechanism/contract"
	"uwwolf/internal/app/game/logic/types"
)

type twoSister struct {
	*role
}

func NewTwoSister(world contract.World, playerID types.PlayerID) (contract.Role, error) {
	return &twoSister{
		role: &role{
			id:           declare.TwoSistersRoleID,
			factionID:    declare.VillagerFactionID,
			phaseID:      declare.NightPhaseID,
			beginRoundID: declare.FirstRound,
			turnID:       declare.TwoSistersTurnID,
			world:        world,
			playerID:     playerID,
			abilities: []*ability{
				{
					action: action.NewRoleIdentify(
						world,
						declare.TwoSistersRoleID,
					),
					activeLimit: declare.Once,
				},
			},
		},
	}, nil
}
