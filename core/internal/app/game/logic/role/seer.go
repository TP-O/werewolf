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

func NewSeer(moderator contract.Moderator, playerId types.PlayerId) (contract.Role, error) {
	return &seer{
		role: &role{
			id:           constants.SeerRoleId,
			factionID:    constants.VillagerFactionId,
			phaseID:      constants.NightPhaseId,
			beginRoundID: constants.SecondRound,
			turnID:       constants.SeerTurnID,
			moderator:    moderator,
			playerId:     playerId,
			abilities: []*ability{
				{
					action: action.NewFactionPredict(
						moderator.World(),
						constants.WerewolfFactionId,
					),
					activeLimit: constants.UnlimitedTimes,
				},
			},
		},
	}, nil
}
