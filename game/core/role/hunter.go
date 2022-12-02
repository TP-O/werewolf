package role

import (
	"uwwolf/game/contract"
	"uwwolf/game/core/action"
	"uwwolf/game/enum"
	"uwwolf/game/types"
)

type hunter struct {
	*role
}

func NewHunter(game contract.Game, playerID enum.PlayerID) (contract.Role, error) {
	return &hunter{
		&role{
			id:         enum.HunterRoleID,
			phaseID:    enum.DayPhaseID,
			factionID:  enum.VillagerFactionID,
			beginRound: enum.FirstRound,
			priority:   enum.HunterTurnPriority,
			game:       game,
			player:     game.Player(playerID),
			abilities: map[enum.ActionID]*ability{
				enum.KillActionID: {
					action:      action.NewKill(game),
					activeLimit: enum.ReachedLimit,
				},
			},
		},
	}, nil
}

func (h *hunter) AfterDeath() {
	diedAt := h.game.Scheduler().PhaseID()
	turnSetting := &types.TurnSetting{
		PhaseID:    h.phaseID,
		RoleID:     h.id,
		BeginRound: h.beginRound,
	}

	if diedAt == h.phaseID {
		// Hunter can play his turn in next turn if he dies at his phase
		turnSetting.Position = enum.NextPosition
	} else {
		// Hunter can play his turn in the next day's his turn
		// if he dies at a time which is not his phase
		turnSetting.Position = enum.SortedPosition
	}

	h.abilities[enum.KillActionID].activeLimit = enum.OneMore
	h.game.Scheduler().AddTurn(turnSetting)
}
