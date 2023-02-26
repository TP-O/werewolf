package role

import (
	"uwwolf/game/action"
	"uwwolf/game/contract"
	"uwwolf/game/types"
)

func NewVillager(game contract.Game, playerID types.PlayerID) (contract.Role, error) {
	voteAction, err := action.NewVote(game, action.VoteActionSetting{
		FactionID: VillagerFactionID,
		PlayerID:  playerID,
		Weight:    1,
	})
	if err != nil {
		return nil, err
	}

	return &role{
		id:           VillagerRoleID,
		factionID:    VillagerFactionID,
		phaseID:      DayPhaseID,
		beginRoundID: types.RoundID(0),
		turnID:       VillagerTurnID,
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
