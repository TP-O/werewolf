package rule

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

const NumberOfWerewolvesTag = "number_of_werewolves"

func GameInstanceInitValidate(sl validator.StructLevel) {
	// strt := sl.Current().Interface().(types.GameInstance)

	// maxNumberOfWerewolves := int(math.Round(float64(strt.Capacity)/2)) - 2

	// if strt.NumberOfWerewolves > maxNumberOfWerewolves {
	// 	sl.ReportError(
	// 		nil,
	// 		"NumberOfWerewolves",
	// 		"NumberOfWerewolves",
	// 		NumberOfWerewolvesTag,
	// 		strconv.Itoa(strt.NumberOfWerewolves),
	// 	)
	// }
}

func AddNumberOfWerewolvesTag(ut ut.Translator) error {
	message := "{0} is too large for capacity game"

	return ut.Add(GameCapacityTag, message, true)
}

func RegisterNumberOfWerewolvesTagMessage(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(GameCapacityTag, fe.Field())

	return t
}
