package role

import (
	"uwwolf/app/enum"
	"uwwolf/app/game/action"
	"uwwolf/app/game/contract"
	"uwwolf/app/types"
)

func newTwoSister(game contract.Game, playerId types.PlayerId) contract.Role {
	return &role{
		id:           enum.TwoSistersRoleId,
		factionId:    enum.VillagerFactionId,
		phaseId:      enum.NightPhaseId,
		game:         game,
		player:       game.Player(playerId),
		beginRoundId: enum.FirstRound,
		priority:     enum.TwoSistersPriority,
		actions: map[types.ActionId]contract.Action{
			enum.RecognitionActionId: action.NewRecognition(game, enum.TwoSistersRoleId),
		},
	}
}
