package role

import (
	"uwwolf/game/action"
	"uwwolf/game/contract"
	"uwwolf/game/types"
	"uwwolf/game/vars"
)

func NewVillager(game contract.Game, playerID types.PlayerID) (contract.Role, error) {
	voteAction, err := action.NewVote(game, &action.VoteActionSetting{
		FactionID: vars.VillagerFactionID,
		PlayerID:  playerID,
		Weight:    1,
	})
	if err != nil {
		return nil, err
	}

	return &role{
		id:           vars.VillagerRoleID,
		factionID:    vars.VillagerFactionID,
		phaseID:      vars.DayPhaseID,
		beginRoundID: types.RoundID(0),
		turnID:       vars.VillagerTurnID,
		game:         game,
		player:       game.Player(playerID),
		abilities: []ability{
			vars.VoteActionID: {
				action:      voteAction,
				activeLimit: vars.Unlimited,
			},
		},
	}, nil
}
