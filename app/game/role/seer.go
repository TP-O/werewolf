package role

import (
	"uwwolf/app/enum"
	"uwwolf/app/game/action"
	"uwwolf/app/game/contract"
	"uwwolf/app/types"
)

func newSeer(game contract.Game, playerId types.PlayerId) contract.Role {
	return &role{
		id:           enum.SeerRoleId,
		factionId:    enum.VillagerFactionId,
		phaseId:      enum.NightPhaseId,
		game:         game,
		player:       game.Player(playerId),
		beginRoundId: types.RoundId(2),
		priority:     enum.SeerPriority,
		actions: map[types.ActionId]contract.Action{
			enum.ProphecyActionId: action.NewProphecy(game),
		},
	}
}
