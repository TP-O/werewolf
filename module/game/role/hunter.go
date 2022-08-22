package role

import (
	"uwwolf/module/game/action"
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

const HunterRoleName = "Hunter"

func NewHunterRole(game contract.Game, setting *types.RoleSetting) contract.Role {
	return &role{
		id:        types.HunterRole,
		factionId: setting.FactionId,
		phaseId:   types.DayPhase,
		name:      HunterRoleName,
		game:      game,
		player:    game.GetPlayer(setting.OwnerId),
		skill: &skill{
			action:       action.NewShooting(game),
			beginRoundId: setting.BeginRound,
			expiration:   setting.Expiration,
		},
	}
}
