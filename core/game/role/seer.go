package role

import (
	"uwwolf/game/action"
	"uwwolf/game/contract"
	"uwwolf/game/types"
	"uwwolf/game/vars"
)

func NewSeer(game contract.Game, playerID types.PlayerID) (contract.Role, error) {
	return &role{
		id:           vars.SeerRoleID,
		factionID:    vars.VillagerFactionID,
		phaseID:      vars.NightPhaseID,
		beginRoundID: types.RoundID(1),
		turnID:       vars.SeerTurnID,
		game:         game,
		player:       game.Player(playerID),
		abilities: []ability{
			{
				action: action.NewFactionPredict(
					game,
					vars.WerewolfFactionID,
				),
				activeLimit: vars.Unlimited,
			},
		},
	}, nil
}
