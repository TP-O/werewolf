package rule

import (
	"math"
	"reflect"
	"uwwolf/game/enum"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

const NumberOfWerewolvesTag = "number_of_werewolves"

func ValidateNumberOfWerewolves(fl validator.FieldLevel) bool {
	var playerIds reflect.Value

	if fl.Parent().Kind() == reflect.Ptr {
		playerIds = fl.Parent().Elem().FieldByName(fl.Param())
	} else {
		playerIds = fl.Parent().FieldByName(fl.Param())
	}

	capacity := len(playerIds.Interface().([]enum.PlayerID))
	maxNumberOfWerewolves := int64(math.Round(float64(capacity)/2)) - 2

	if fl.Field().Int() > maxNumberOfWerewolves {
		return false
	}

	return true
}

func AddNumberOfWerewolvesTag(ut ut.Translator) error {
	message := "{0} is not balanced"

	return ut.Add(NumberOfWerewolvesTag, message, true)
}

func RegisterNumberOfWerewolvesTagMessage(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(NumberOfWerewolvesTag, fe.Field())

	return t
}
