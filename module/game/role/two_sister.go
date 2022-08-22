package role

import (
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

const TwoSistersRoleName = "Two Sisters"

type twoSister struct {
	role
}

func NewTwoSisterRole(game contract.Game, setting *types.RoleSetting) contract.Role {
	return &role{
		id:        types.TwoSistersRole,
		factionId: setting.FactionId,
		phaseId:   types.NightPhase,
		name:      TwoSistersRoleName,
		game:      game,
		player:    game.GetPlayer(setting.OwnerId),
		skill: &skill{
			action:       nil,
			beginRoundId: setting.BeginRound,
			expiration:   setting.Expiration,
		},
	}
}
