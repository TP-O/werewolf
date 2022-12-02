package role

import (
	"uwwolf/game/contract"
	"uwwolf/game/core/action"
	"uwwolf/game/enum"
	"uwwolf/game/types"
)

func NewWerewolf(game contract.Game, playerID enum.PlayerID) (contract.Role, error) {
	voteAction, err := action.NewVote(game, &types.VoteActionSetting{
		FactionID: enum.WerewolfFactionID,
		PlayerID:  playerID,
		Weight:    1,
	})

	if err != nil {
		return nil, err
	}

	return &role{
		id:         enum.WerewolfRoleID,
		factionID:  enum.WerewolfFactionID,
		phaseID:    enum.NightPhaseID,
		beginRound: enum.FirstRound,
		priority:   enum.WerewolfTurnPriority,
		game:       game,
		player:     game.Player(playerID),
		abilities: map[enum.ActionID]*ability{
			enum.VoteActionID: {
				action:      voteAction,
				activeLimit: enum.Unlimited,
			},
		},
	}, nil
}
