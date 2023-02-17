package tag

import (
	"strconv"
	"uwwolf/config"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

const TurnDurationTag = "turn_duration"

func AddTurnDurationTag(ut ut.Translator) error {
	message := "{0} must be from " +
		strconv.Itoa(int(config.Game().MinTurnDuration)) +
		" to " +
		strconv.Itoa(int(config.Game().MaxTurnDuration)) +
		" seconds"

	return ut.Add(TurnDurationTag, message, true)
}

func RegisterTurnDurationMessage(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(TurnDurationTag, fe.Field())

	return t
}
