package role

import (
	"uwwolf/game/action"
	"uwwolf/game/contract"
	"uwwolf/game/types"
	"uwwolf/game/vars"
)

type twoSister struct {
	*role
}

func NewTwoSister(game contract.Game, playerID types.PlayerID) (contract.Role, error) {
	return &twoSister{
		role: &role{
			id:           vars.TwoSistersRoleID,
			factionID:    vars.VillagerFactionID,
			phaseID:      vars.NightPhaseID,
			beginRoundID: vars.FirstRound,
			turnID:       vars.TwoSistersTurnID,
			game:         game,
			player:       game.Player(playerID),
			abilities: []*ability{
				{
					action: action.NewRoleIdentify(
						game,
						vars.TwoSistersRoleID,
					),
					activeLimit: vars.One,
				},
			},
		},
	}, nil
}

// RegisterTurn adds role's turn to the game schedule.
func (ts twoSister) RegisterTurn() {
	ts.game.Scheduler().AddPlayerTurn(&types.NewPlayerTurn{
		PhaseID:      ts.phaseID,
		TurnID:       ts.turnID,
		BeginRoundID: ts.beginRoundID,
		PlayerID:     ts.player.ID(),
		RoleID:       ts.id,
		ExpiredAfter: vars.One,
	})
}
