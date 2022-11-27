package role

import (
	"uwwolf/app/game/config"
	"uwwolf/app/game/contract"
	"uwwolf/app/game/core/action"
	"uwwolf/app/game/types"
)

func NewWerewolf(game contract.Game, playerID types.PlayerID) contract.Role {
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
				action: action.NewVote(game, &types.VoteActionSetting{
					FactionID: config.WerewolfFactionID,
					PlayerID:  playerID,
					Weight:    1,
				}),
				activeLimit: config.Unlimited,
			},
		},
	}
}
