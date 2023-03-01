package role

import (
	"uwwolf/game/action"
	"uwwolf/game/contract"
	"uwwolf/game/types"
)

type hunter struct {
	role
}

func NewHunter(game contract.Game, playerID types.PlayerID) (contract.Role, error) {
	return &hunter{
			role{
				id:           HunterRoleID,
				phaseID:      DayPhaseID,
				factionID:    VillagerFactionID,
				beginRoundID: types.RoundID(0),
				turnID:       HunterTurnID,
				game:         game,
				player:       game.Player(playerID),
				abilities: []ability{
					{
						action:      action.NewKill(game),
						activeLimit: ReachedLimit,
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
		playerTurn.TurnID = HunterTurnID
	}

	h.abilities[action.KillActionID].activeLimit = One
	h.game.Scheduler().AddPlayerTurn(playerTurn)
}

func (h *hunter) AfterSaved() {
	// Undo `AfterDeath`
	h.abilities[action.KillActionID].activeLimit = ReachedLimit
	h.game.Scheduler().RemovePlayerTurn(types.RemovedPlayerTurn{
		PhaseID:  h.phaseID,
		TurnID:   types.TurnID(-1),
		PlayerID: h.player.ID(),
		RoleID:   h.id,
	})
}
