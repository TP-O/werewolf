package role

import (
	"uwwolf/app/game/action"
	"uwwolf/app/game/contract"
	"uwwolf/app/types"
)

type hunerRole struct {
	role
}

func NewHunterRole(game contract.Game, setting *types.RoleSetting) contract.Role {
	return &hunerRole{
		role: role{
			id:        setting.Id,
			factionId: setting.FactionId,
			phaseId:   setting.PhaseId,
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
	turnSetting := &types.TurnSetting{
		PhaseId:    h.role.phaseId,
		RoleId:     h.role.id,
		PlayerIds:  []types.PlayerId{h.role.player.Id()},
		BeginRound: h.role.skill.beginRoundId,
		Expiration: h.role.skill.expiration,
	}

	if diedAt == h.role.phaseId {
		// Hunter can play his turn in next turn if he dies at his phase
		turnSetting.Position = types.NextPosition
	} else {
		// Hunter can play his turn in the next day's first turn
		// if he dies at a time which is not his phase
		turnSetting.Position = types.FirstPosition
	}

	h.role.game.Round().AddTurn(turnSetting)
}
