package tag

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

const NumberWerewolvesTag = "number_werewolves"

func AddNumberWerewolvesTag(ut ut.Translator) error {
	message := "{0} is not balanced"

	return ut.Add(NumberWerewolvesTag, message, true)
}

func RegisterNumberWerewolvesMessage(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(NumberWerewolvesTag, fe.Field())

	return t
}
