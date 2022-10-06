package service

import (
	"uwwolf/app/game/contract"
	"uwwolf/app/game/core"
	"uwwolf/app/model"
	"uwwolf/app/types"
	"uwwolf/db"
)

var gameManger = core.NewManager()

func CreateGame(setting *types.GameSetting) contract.Game {
	game := &model.Game{}
	db.Client().Omit("WinningFactionId").Create(game)

	return gameManger.AddGame(game.Id, setting)
}
