package role

import (
	"uwwolf/module/game/action"
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

const HunterRoleName = "Hunter"

type hunerRole struct {
	role
}

func NewHunterRole(game contract.Game, setting *types.RoleSetting) contract.Role {
	return &hunerRole{
		role: role{
			id:        types.HunterRole,
			factionId: setting.FactionId,
			phaseId:   types.DayPhase,
			name:      HunterRoleName,
			game:      game,
			player:    game.Player(setting.OwnerId),
			skill: &skill{
				action:       action.NewShooting(game),
				beginRoundId: setting.BeginRound,
				expiration:   setting.Expiration,
			},
		},
	}
}

func (h *hunerRole) AfterDeath() {
	diedAt := h.role.game.Round().CurrentPhaseId()

	if diedAt == types.NightPhase {
		// Hunter can play his turn in next day's first turn if he dies at night
		h.role.game.Round().AddTurn(&types.TurnSetting{
			PhaseId:    h.role.phaseId,
			RoleId:     h.role.id,
			PlayerIds:  []types.PlayerId{h.role.player.Id()},
			BeginRound: h.role.skill.beginRoundId,
			Position:   types.FirstPosition,
			Expiration: h.role.skill.expiration,
		})
	} else if diedAt == h.role.phaseId {
		// Hunter can play his turn in next turn if he dies at his phase
		h.role.game.Round().AddTurn(&types.TurnSetting{
			RoleId:     h.role.id,
			PlayerIds:  []types.PlayerId{h.role.player.Id()},
			BeginRound: h.role.skill.beginRoundId,
			Position:   types.NextPosition,
			Expiration: h.role.skill.expiration,
		})
	}
}
