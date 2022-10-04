package core

import (
	validate "github.com/go-playground/validator/v10"

	"uwwolf/app/game/contract"
	"uwwolf/app/model"
	"uwwolf/app/types"
	"uwwolf/app/validator"
	"uwwolf/db"
)

type manager struct {
	games map[types.GameId]contract.Game
}

var mangerInstance *manager

func NewManager() contract.GameManger {
	if mangerInstance == nil {
		mangerInstance = &manager{
			games: make(map[types.GameId]contract.Game),
		}
	}

	return mangerInstance
}

func (m *manager) Game(gameId types.GameId) contract.Game {
	return m.games[gameId]
}

func (m *manager) AddGame(setting *types.GameSetting) validate.ValidationErrorsTranslations {
	if err := validator.ValidateStruct(setting); err != nil {
		return err
	}

	game := &model.Game{}
	db.Client().Omit("WinningFactionId").Create(game)
	m.games[game.Id] = NewGame(game.Id, setting)

	return nil
}
