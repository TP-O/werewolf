package role

import (
	"uwwolf/app/game/config"
	"uwwolf/app/game/contract"
	"uwwolf/app/game/core/action"
	"uwwolf/app/game/types"
)

func NewWerewolf(game contract.Game, playerID types.PlayerID) (contract.Role, error) {
	voteAction, err := action.NewVote(game, &types.VoteActionSetting{
		FactionID: config.WerewolfFactionID,
		PlayerID:  playerID,
		Weight:    1,
	})

	if err != nil {
		return nil, err
	}

	return &role{
		id:         config.WerewolfRoleID,
		factionID:  config.WerewolfFactionID,
		phaseID:    config.NightPhaseID,
		beginRound: config.FirstRound,
		priority:   config.WerewolfTurnPriority,
		game:       game,
		player:     game.Player(playerID),
		abilities: map[types.ActionID]*ability{
			config.VoteActionID: {
				action:      voteAction,
				activeLimit: config.Unlimited,
			},
		},
	}, nil
}
