package role

import (
	"uwwolf/app/game/config"
	"uwwolf/app/game/contract"
	"uwwolf/app/game/core/action"
	"uwwolf/app/game/types"
)

func NewVillager(game contract.Game, playerID types.PlayerID) contract.Role {
	return &role{
		id:         config.VillagerRoleID,
		factionID:  config.VillagerFactionID,
		phaseID:    config.DayPhaseID,
		beginRound: config.FirstRound,
		priority:   config.VillagerTurnPriority,
		game:       game,
		player:     game.Player(playerID),
		abilities: map[types.ActionID]*ability{
			config.VoteActionID: {
				action: action.NewVote(game, &types.VoteActionSetting{
					FactionID: config.VillagerFactionID,
					PlayerID:  playerID,
					Weight:    1,
				}),
				activeLimit: config.Unlimited,
			},
		},
	}
}
