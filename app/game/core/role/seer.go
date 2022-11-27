package role

import (
	"uwwolf/app/game/config"
	"uwwolf/app/game/contract"
	"uwwolf/app/game/core/action"
	"uwwolf/app/game/types"
)

func NewSeer(game contract.Game, playerID types.PlayerID) contract.Role {
	return &role{
		id:         config.SeerRoleID,
		factionID:  config.VillagerFactionID,
		phaseID:    config.NightPhaseID,
		beginRound: types.Round(2),
		priority:   config.SeerTurnPriority,
		game:       game,
		player:     game.Player(playerID),
		abilities: map[types.ActionID]*ability{
			config.PredictActionID: {
				action: action.NewFactionPredict(
					game,
					config.WerewolfFactionID,
				),
				activeLimit: config.Unlimited,
			},
		},
	}
}
