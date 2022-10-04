package role

import (
	"uwwolf/app/enum"
	"uwwolf/app/game/action"
	"uwwolf/app/game/contract"
	"uwwolf/app/types"
)

func newWerewolf(game contract.Game, playerId types.PlayerId) contract.Role {
	return &role{
		id:           enum.WerewolfRoleId,
		factionId:    enum.WerewolfFactionId,
		phaseId:      enum.NightPhaseId,
		game:         game,
		player:       game.Player(playerId),
		beginRoundId: enum.FirstRound,
		priority:     enum.WerewolfPriority,
		actions: map[types.ActionId]contract.Action{
			enum.VoteActionId: action.NewVote(game, &types.VoteActionSetting{
				FactionId: enum.WerewolfFactionId,
				PlayerId:  playerId,
				Weight:    1,
			}),
		},
	}
}
