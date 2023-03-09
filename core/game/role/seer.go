package role

import (
	"uwwolf/game/action"
	"uwwolf/game/contract"
	"uwwolf/game/types"
	"uwwolf/game/vars"
)

type seer struct {
	*role
}

func NewSeer(game contract.Game, playerID types.PlayerID) (contract.Role, error) {
	return &seer{
		role: &role{
			id:           vars.SeerRoleID,
			factionID:    vars.VillagerFactionID,
			phaseID:      vars.NightPhaseID,
			beginRoundID: vars.SecondRound,
			turnID:       vars.SeerTurnID,
			game:         game,
			player:       game.Player(playerID),
			abilities: []*ability{
				{
					action: action.NewFactionPredict(
						game,
						vars.WerewolfFactionID,
					),
					activeLimit: vars.UnlimitedTimes,
				},
			},
		},
	}, nil
}
