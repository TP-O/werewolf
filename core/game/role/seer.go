package role

import (
	"uwwolf/game/action"
	"uwwolf/game/contract"
	"uwwolf/game/types"
)

func NewSeer(game contract.Game, playerID types.PlayerID) (contract.Role, error) {
	return &role{
		id:           SeerRoleID,
		factionID:    VillagerFactionID,
		phaseID:      NightPhaseID,
		beginRoundID: types.RoundID(1),
		turnID:       SeerTurnID,
		game:         game,
		player:       game.Player(playerID),
		abilities: []ability{
			{
				action: action.NewFactionPredict(
					game,
					WerewolfFactionID,
				),
				activeLimit: Unlimited,
			},
		},
	}, nil
}
