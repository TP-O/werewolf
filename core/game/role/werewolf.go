package role

import (
	"uwwolf/game/action"
	"uwwolf/game/contract"
	"uwwolf/game/types"
)

func NewWerewolf(game contract.Game, playerID types.PlayerID) (contract.Role, error) {
	voteAction, err := action.NewVote(game, action.VoteActionSetting{
		FactionID: WerewolfFactionID,
		PlayerID:  playerID,
		Weight:    1,
	})
	if err != nil {
		return nil, err
	}

	return &role{
		id:           WerewolfRoleID,
		factionID:    WerewolfFactionID,
		phaseID:      NightPhaseID,
		beginRoundID: types.RoundID(0),
		turnID:       WerewolfTurnID,
		game:         game,
		player:       game.Player(playerID),
		abilities: []ability{
			action.VoteActionID: {
				action:      voteAction,
				activeLimit: Unlimited,
			},
		},
	}, nil
}
