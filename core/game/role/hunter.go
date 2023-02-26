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
	turnSetting := types.TurnSetting{
		PhaseID:   h.phaseID,
		PlayerIDs: []types.PlayerID{h.player.ID()},
	}

	if diedAtPhaseID == h.phaseID {
		// Hunter can play in next turn if he dies at his phase
		turnSetting.BeginRoundID = h.game.Scheduler().RoundID()
		turnSetting.TurnID = h.game.Scheduler().Turn().ID + 1
	} else {
		// Hunter can play in his turn of the next day
		// if he dies at a time which is not his phase
		turnSetting.BeginRoundID = h.game.Scheduler().RoundID() + 1
		turnSetting.TurnID = HunterTurnID
	}

	h.abilities[action.KillActionID].activeLimit = One
	h.game.Scheduler().AddTurn(turnSetting)
}
