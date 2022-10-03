package role

import (
	"uwwolf/app/enum"
	"uwwolf/app/game/action"
	"uwwolf/app/game/contract"
	"uwwolf/app/types"
)

type hunerRole struct {
	role
}

func newHunter(game contract.Game, playerId types.PlayerId) contract.Role {
	return &hunerRole{
		role: role{
			id:           enum.HunterRoleId,
			factionId:    enum.VillagerFactionId,
			phaseId:      enum.DayPhaseId,
			game:         game,
			player:       game.Player(playerId),
			beginRoundId: enum.FirstRound,
			priority:     0,
			score:        1,
			set:          1,
			actions: map[uint]contract.Action{
				enum.ShootingActionId: action.NewShooting(game),
			},
		},
	}
}

func (hr *hunerRole) AfterDeath() {
	diedAt := hr.role.game.Round().CurrentPhaseId()
	turnSetting := &types.TurnSetting{
		PhaseId:    hr.role.phaseId,
		RoleId:     hr.role.id,
		PlayerIds:  []types.PlayerId{hr.role.player.Id()},
		BeginRound: hr.beginRoundId,
		Expiration: hr.actions[enum.ShootingActionId].Expiration(),
	}

	if diedAt == hr.role.phaseId {
		// Hunter can play his turn in next turn if he dies at his phase
		turnSetting.Position = enum.NextPosition
	} else {
		// Hunter can play his turn in the next day's first turn
		// if he dies at a time which is not his phase
		turnSetting.Position = enum.FirstPosition
	}

	hr.role.game.Round().AddTurn(turnSetting)
}
