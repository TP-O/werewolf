package tag

import (
	"strconv"
	"uwwolf/util"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

const GameCapacityTag = "game_capacity"

func AddGameCapacityTag(ut ut.Translator) error {
	message := "{0} must have " +
		strconv.Itoa(int(util.Config().Game.MinCapacity)) +
		" to " +
		strconv.Itoa(int(util.Config().Game.MaxCapacity)) +
		" players"

	return ut.Add(GameCapacityTag, message, true)
}

func RegisterGameCapacityMessage(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(GameCapacityTag, fe.Field())

	return t
}
