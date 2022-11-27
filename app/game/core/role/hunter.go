package role

import (
	"uwwolf/app/game/config"
	"uwwolf/app/game/contract"
	"uwwolf/app/game/core/action"
	"uwwolf/app/game/types"
)

type hunter struct {
	*role
}

func NewHunter(game contract.Game, playerID types.PlayerID) contract.Role {
	return &hunter{
		&role{
			id:         config.HunterRoleID,
			phaseID:    config.DayPhaseID,
			factionID:  config.WerewolfFactionID,
			beginRound: config.FirstRound,
			priority:   config.HunterTurnPriority,
			game:       game,
			player:     game.Player(playerID),
			abilities: map[types.ActionID]*ability{
				config.KillActionID: {
					action:      action.NewKill(game),
					activeLimit: config.ReachedLimit,
				},
			},
		},
	}
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
		turnSetting.Position = config.NextPosition
	} else {
		// Hunter can play his turn in the next day's his turn
		// if he dies at a time which is not his phase
		turnSetting.Position = config.SortedPosition
	}

	h.abilities[config.KillActionID].activeLimit = config.OneMore
	h.game.Scheduler().AddTurn(turnSetting)
}
