package role

import (
	"uwwolf/app/game/config"
	"uwwolf/app/game/contract"
	"uwwolf/app/game/core/action"
	"uwwolf/app/game/types"
)

func NewTwoSister(game contract.Game, playerID types.PlayerID) contract.Role {
	return &role{
		id:         config.TwoSistersRoleID,
		factionID:  config.VillagerFactionID,
		phaseID:    config.NightPhaseID,
		beginRound: config.FirstRound,
		priority:   config.TwoSistersTurnPriority,
		game:       game,
		player:     game.Player(playerID),
		abilities: map[types.ActionID]*ability{
			config.RecognizeActionID: {
				action: action.NewRoleRecognize(
					game,
					config.TwoSistersRoleID,
				),
				activeLimit: config.OneMore,
			},
		},
	}
}
