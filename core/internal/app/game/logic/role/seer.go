package role

import (
	"uwwolf/internal/app/game/logic/action"
	"uwwolf/internal/app/game/logic/constants"
	"uwwolf/internal/app/game/logic/contract"
	"uwwolf/internal/app/game/logic/types"
)

type seer struct {
	*role
}

func NewSeer(world contract.World, playerID types.PlayerID) (contract.Role, error) {
	return &seer{
		role: &role{
			id:           constants.SeerRoleID,
			factionID:    constants.VillagerFactionID,
			phaseID:      constants.NightPhaseID,
			beginRoundID: constants.SecondRound,
			turnID:       constants.SeerTurnID,
			world:        world,
			playerID:     playerID,
			abilities: []*ability{
				{
					action: action.NewFactionPredict(
						world,
						constants.WerewolfFactionID,
					),
					activeLimit: constants.UnlimitedTimes,
				},
			},
		},
	}, nil
}
