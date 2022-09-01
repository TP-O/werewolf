package role

import (
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

const TwoSistersRoleName = "Two Sisters"

func NewTwoSisterRole(game contract.Game, setting *types.RoleSetting) contract.Role {
	return &role{
		id:        types.TwoSistersRole,
		factionId: setting.FactionId,
		phaseId:   types.NightPhase,
		name:      TwoSistersRoleName,
		game:      game,
		player:    game.Player(setting.OwnerId),
		skill: &skill{
			beginRoundId: setting.BeginRound,
			expiration:   setting.Expiration,
		},
	}
}
