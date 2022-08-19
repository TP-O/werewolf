package role

import (
	"uwwolf/module/game/action"
	"uwwolf/module/game/contract"
	"uwwolf/types"
)

const SeerRoleName = "Seer"

func NewSeerRole(game contract.Game, playerId types.PlayerId) contract.Role {
	return &role{
		id:      types.SeerRole,
		phaseId: types.NightPhase,
		name:    SeerRoleName,
		game:    game,
		player:  game.GetPlayer(playerId),
		skill: &skill{
			action:       action.NewProphecy(game),
			numberOfUses: types.UnlimitedTimes,
			beginRoundId: types.FirstRound,
		},
	}
}
