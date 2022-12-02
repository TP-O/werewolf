package role

import (
	"uwwolf/game/contract"
	"uwwolf/game/core/action"
	"uwwolf/game/enum"
)

func NewTwoSister(game contract.Game, playerID enum.PlayerID) (contract.Role, error) {
	return &role{
		id:         enum.TwoSistersRoleID,
		factionID:  enum.VillagerFactionID,
		phaseID:    enum.NightPhaseID,
		beginRound: enum.FirstRound,
		priority:   enum.TwoSistersTurnPriority,
		game:       game,
		player:     game.Player(playerID),
		abilities: map[enum.ActionID]*ability{
			enum.RecognizeActionID: {
				action: action.NewRoleRecognize(
					game,
					enum.TwoSistersRoleID,
				),
				activeLimit: enum.OneMore,
			},
		},
	}, nil
}
