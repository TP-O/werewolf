package service

import (
	"github.com/go-playground/validator/v10"

	"uwwolf/app/game/core"
	"uwwolf/app/types"
)

var gameManger = core.NewManager()

func CreateGame(setting *types.GameSetting) validator.ValidationErrorsTranslations {
	err := gameManger.AddGame(setting)

	return err
}
