package action

import (
	"uwwolf/module/game"
	"uwwolf/types"
)

const ShootingAction = "Shooting"

func NewShootingAction() Action {
	return &action[int]{
		name:     ShootingAction,
		validate: validateShooting,
		execute:  executeShooting,
		skip:     skipShooting,
	}
}

func validateShooting(data *types.ActionData) bool {
	return data.Skipped ||
		(!data.Skipped && len(data.Targets) == 1)
}

func executeShooting(_ game.IGame, data *types.ActionData, _ int) bool {
	// fmt.Println(data.Actor + " shot " + data.Targets[0])

	return true
}

func skipShooting(_ game.IGame, data *types.ActionData, _ int) bool {
	// fmt.Println(data.Actor + " skipped")

	return true
}
