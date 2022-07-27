package action

import (
	"uwwolf/module/game"
	"uwwolf/types"
)

const ProphecyAction = "Prophecy"

func NewProphecyAction() Action {
	return &action[int]{
		name:     ProphecyAction,
		validate: validateProphecy,
		execute:  executeProphecy,
		skip:     skipProphecy,
	}
}

func validateProphecy(data *types.ActionData) bool {
	return data.Skipped ||
		(!data.Skipped && len(data.Targets) == 1)
}

func executeProphecy(_ game.IGame, data *types.ActionData, _ int) bool {
	// fmt.Println(data.Actor + " pophesied " + data.Targets[0] + " is werewolf")

	return true
}

func skipProphecy(_ game.IGame, data *types.ActionData, _ int) bool {
	// fmt.Println(data.Actor + " skipped")

	return true
}
