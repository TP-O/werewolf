package role

import (
	"uwwolf/internal/app/game/logic/declare"
	"uwwolf/internal/app/game/logic/mechanism/action"
	"uwwolf/internal/app/game/logic/mechanism/contract"
	"uwwolf/internal/app/game/logic/types"
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
