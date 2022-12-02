package rule

import (
	"strconv"
	"uwwolf/config"
	"uwwolf/game/enum"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

const CapacityTag = "capacity"

func ValidateCapacity(fl validator.FieldLevel) bool {
	capacity := len(fl.Field().Interface().([]enum.PlayerID))

	if capacity > config.Game.MaxCapacity ||
		capacity < config.Game.MinCapacity {

		return false
	}

	return true
}

func AddCapacityTag(ut ut.Translator) error {
	message := "{0} must have length from " +
		strconv.Itoa(config.Game.MinCapacity) +
		" to " +
		strconv.Itoa(config.Game.MaxCapacity)

	return ut.Add(CapacityTag, message, true)
}

func RegisterCapacityMessage(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(CapacityTag, fe.Field())

	return t
}
