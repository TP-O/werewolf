package role

import (
	"uwwolf/game/action"
	"uwwolf/game/contract"
	"uwwolf/game/types"
	"uwwolf/game/vars"
)

type hunter struct {
	role
}

func NewHunter(game contract.Game, playerID types.PlayerID) (contract.Role, error) {
	return &hunter{
			role{
				id:           vars.HunterRoleID,
				phaseID:      vars.DayPhaseID,
				factionID:    vars.VillagerFactionID,
				beginRoundID: types.RoundID(0),
				turnID:       vars.HunterTurnID,
				game:         game,
				player:       game.Player(playerID),
				abilities: []ability{
					{
						action:      action.NewKill(game),
						activeLimit: vars.ReachedLimit,
					},
				},
			},
		},
		nil
}

func (h *hunter) AfterDeath() {
	diedAtPhaseID := h.game.Scheduler().PhaseID()
	playerTurn := types.NewPlayerTurn{
		PhaseID:  h.phaseID,
		PlayerID: h.player.ID(),
		RoleID:   h.id,
	}
	if diedAtPhaseID == h.phaseID {
		// Hunter can play in next turn if he dies at his phase
		playerTurn.BeginRoundID = h.game.Scheduler().RoundID()
		playerTurn.TurnID = h.game.Scheduler().TurnID() + 1
	} else {
		// Hunter can play in his turn of the next day
		// if he dies at a time which is not his phase
		playerTurn.BeginRoundID = h.game.Scheduler().RoundID() + 1
		playerTurn.TurnID = vars.HunterTurnID
	}

	h.abilities[vars.KillActionID].activeLimit = vars.One
	h.game.Scheduler().AddPlayerTurn(playerTurn)
}

func (h *hunter) AfterSaved() {
	// Undo `AfterDeath`
	h.abilities[vars.KillActionID].activeLimit = vars.ReachedLimit
	h.game.Scheduler().RemovePlayerTurn(types.RemovedPlayerTurn{
		PhaseID:  h.phaseID,
		TurnID:   types.TurnID(-1),
		PlayerID: h.player.ID(),
		RoleID:   h.id,
	})
}
