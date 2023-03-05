package role

import (
	"uwwolf/game/action"
	"uwwolf/game/contract"
	"uwwolf/game/types"
	"uwwolf/game/vars"
)

func NewTwoSister(game contract.Game, playerID types.PlayerID) (contract.Role, error) {
	return &role{
		id:           vars.TwoSistersRoleID,
		factionID:    vars.VillagerFactionID,
		phaseID:      vars.NightPhaseID,
		beginRoundID: types.RoundID(0),
		turnID:       vars.TwoSistersTurnID,
		game:         game,
		player:       game.Player(playerID),
		abilities: []ability{
			{
				action: action.NewRoleIdentify(
					game,
					vars.TwoSistersRoleID,
				),
				activeLimit: vars.One,
			},
		},
	}, nil
}
