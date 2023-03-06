package role

import (
	"uwwolf/game/action"
	"uwwolf/game/contract"
	"uwwolf/game/types"
	"uwwolf/game/vars"
)

type werewolf struct {
	*role
}

func NewWerewolf(game contract.Game, playerID types.PlayerID) (contract.Role, error) {
	voteAction, err := action.NewVote(game, &action.VoteActionSetting{
		FactionID: vars.WerewolfFactionID,
		PlayerID:  playerID,
		Weight:    1,
	})
	if err != nil {
		return nil, err
	}

	return &werewolf{
		role: &role{
			id:           vars.WerewolfRoleID,
			factionID:    vars.WerewolfFactionID,
			phaseID:      vars.NightPhaseID,
			beginRoundID: vars.FirstRound,
			turnID:       vars.WerewolfTurnID,
			game:         game,
			player:       game.Player(playerID),
			abilities: []*ability{
				{
					action:      voteAction,
					activeLimit: vars.Unlimited,
				},
			},
		},
	}, nil
}

// RegisterTurn adds role's turn to the game schedule.
func (w werewolf) RegisterTurn() {
	w.game.Scheduler().AddPlayerTurn(&types.NewPlayerTurn{
		PhaseID:      w.phaseID,
		TurnID:       w.turnID,
		BeginRoundID: w.beginRoundID,
		PlayerID:     w.player.ID(),
		RoleID:       w.id,
		ExpiredAfter: vars.Unlimited,
	})
}
