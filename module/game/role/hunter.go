package role

import (
	"uwwolf/module/game/action"
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

const HunterRoleName = "Hunter"

func NewHunterRole(game contract.Game, playerId types.PlayerId) contract.Role {
	return &role{
		id:      types.HunterRole,
		phaseId: types.DayPhase,
		name:    HunterRoleName,
		game:    game,
		player:  game.GetPlayer(playerId),
		skill: &skill{
			action:       action.NewShooting(game),
			numberOfUses: types.OneTimes,
			beginRoundId: 1,
		},
	}
}
