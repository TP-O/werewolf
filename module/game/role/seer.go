package role

import (
	"uwwolf/module/game/action"
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

const SeerRoleName = "Seer"

func NewSeerRole(game contract.Game, setting *types.RoleSetting) contract.Role {
	return &role{
		id:      types.SeerRole,
		phaseId: types.NightPhase,
		name:    SeerRoleName,
		game:    game,
		player:  game.GetPlayer(setting.OwnerId),
		skill: &skill{
			action:       action.NewProphecy(game),
			beginRoundId: setting.BeginRound,
			expiration:   setting.Expiration,
		},
	}
}
