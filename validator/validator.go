package validator

import (
	"fmt"

	"uwwolf/config"
)

func handleError(errs error, isLogged bool) bool {
	if errs != nil {
		if config.App.Debug && isLogged {
			fmt.Println(errs)
		}

		return false
	}

	return true
}

func ValidateVar[T any](input T, tag string) bool {
	return handleError(validator.Var(input, tag), false)
}

func ValidateStruct[T any](input T) bool {
	return handleError(validator.Struct(input), true)
}
