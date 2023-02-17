package role

import (
	"uwwolf/game/contract"
	"uwwolf/game/core/action"
	"uwwolf/game/enum"
)

func NewSeer(game contract.Game, playerID enum.PlayerID) (contract.Role, error) {
	return &role{
		id:         enum.SeerRoleID,
		factionID:  enum.VillagerFactionID,
		phaseID:    enum.NightPhaseID,
		beginRound: enum.Round(2),
		priority:   enum.SeerTurnPriority,
		game:       game,
		player:     game.Player(playerID),
		abilities: map[enum.ActionID]*ability{
			enum.PredictActionID: {
				action: action.NewFactionPredict(
					game,
					enum.WerewolfFactionID,
				),
				activeLimit: enum.Unlimited,
			},
		},
	}, nil
}
