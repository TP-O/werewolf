package role

import (
	"uwwolf/game/action"
	"uwwolf/game/contract"
	"uwwolf/game/types"
	"uwwolf/game/vars"
)

func NewWerewolf(game contract.Game, playerID types.PlayerID) (contract.Role, error) {
	voteAction, err := action.NewVote(game, &action.VoteActionSetting{
		FactionID: vars.WerewolfFactionID,
		PlayerID:  playerID,
		Weight:    1,
	})
	if err != nil {
		return nil, err
	}

	return &role{
		id:           vars.WerewolfRoleID,
		factionID:    vars.WerewolfFactionID,
		phaseID:      vars.NightPhaseID,
		beginRoundID: types.RoundID(0),
		turnID:       vars.WerewolfTurnID,
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
