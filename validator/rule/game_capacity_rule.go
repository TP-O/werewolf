package rule

import (
	"strconv"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	"uwwolf/config"
)

const GameCapacityTag = "game_capacity"

func GameCapacityValidate(fl validator.FieldLevel) bool {
	return fl.Field().Int() <= int64(config.Game.MaxCapacity) &&
		fl.Field().Int() >= int64(config.Game.MinCapacity)
}

func AddGameCapacityTag(ut ut.Translator) error {
	message := "{0} must be between " +
		strconv.Itoa(config.Game.MinCapacity) +
		" and " +
		strconv.Itoa(config.Game.MaxCapacity)

	return ut.Add(GameCapacityTag, message, true)
}

func RegisterGameCapacityMessage(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(GameCapacityTag, fe.Field())

	return t
}
