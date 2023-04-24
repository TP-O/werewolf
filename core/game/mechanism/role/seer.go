package role

import (
	"uwwolf/game/declare"
	"uwwolf/game/mechanism/action"
	"uwwolf/game/mechanism/contract"
	"uwwolf/game/types"
)

type seer struct {
	*role
}

func NewSeer(world contract.World, playerID types.PlayerID) (contract.Role, error) {
	return &seer{
		role: &role{
			id:           declare.SeerRoleID,
			factionID:    declare.VillagerFactionID,
			phaseID:      declare.NightPhaseID,
			beginRoundID: declare.SecondRound,
			turnID:       declare.SeerTurnID,
			world:        world,
			playerID:     playerID,
			abilities: []*ability{
				{
					action: action.NewFactionPredict(
						world,
						declare.WerewolfFactionID,
					),
					activeLimit: declare.UnlimitedTimes,
				},
			},
		},
	}, nil
}
