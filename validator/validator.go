package validator

import (
	"fmt"
	"os"

	validation "github.com/go-playground/validator/v10"

	"uwwolf/game/contract"
	"uwwolf/validator/rule"
)

var validator *validation.Validate

func init() {
	validator = validation.New()

	registerRules()
}

// Register custom validations
func registerRules() {
	validator.RegisterValidation(rule.GameCapacityRule, rule.GameCapacityValidate)
	validator.RegisterStructValidation(rule.GameInstanceInitValidate, contract.GameInstanceInit{})
}

func handleValidaion(errs error, isLogged bool) bool {
	if errs != nil {
		if os.Getenv("DEBUG") == "true" && isLogged {
			fmt.Println(errs)
		}

		return false
	}

	return true
}

func ValidateVar[T any](input T, tag string) bool {
	return handleValidaion(validator.Var(input, tag), false)
}

func ValidateStruct[T any](input T) bool {
	return handleValidaion(validator.Struct(input), true)
}
