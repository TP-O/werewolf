package role

import (
	"uwwolf/game/action"
	"uwwolf/game/contract"
	"uwwolf/game/types"
)

func NewTwoSister(game contract.Game, playerID types.PlayerID) (contract.Role, error) {
	return &role{
		id:           TwoSistersRoleID,
		factionID:    VillagerFactionID,
		phaseID:      NightPhaseID,
		beginRoundID: types.RoundID(0),
		turnID:       TwoSistersTurnID,
		game:         game,
		player:       game.Player(playerID),
		abilities: []ability{
			{
				action: action.NewRoleIdentify(
					game,
					TwoSistersRoleID,
				),
				activeLimit: One,
			},
		},
	}, nil
}
