package core

import (
	validate "github.com/go-playground/validator/v10"

	"uwwolf/app/game/contract"
	"uwwolf/app/types"
	"uwwolf/app/validator"
)

type manager struct {
	games map[types.GameId]contract.Game
}

var mangerInstance *manager

func New() contract.GameManger {
	if mangerInstance == nil {
		mangerInstance = &manager{}
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

	m.games[setting.Id] = NewGame(setting)

	return nil
}
