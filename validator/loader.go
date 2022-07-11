package validator

import (
	validation "github.com/go-playground/validator/v10"

	"uwwolf/contract/typ"
	"uwwolf/validator/rule"
)

var validator *validation.Validate

func LoadValidator() {
	validator = validation.New()

	registerRules()
}

// Register custom validations
func registerRules() {
	validator.RegisterValidation(rule.GameCapacityRule, rule.GameCapacityValidate)
	validator.RegisterStructValidation(rule.GameInstanceInitValidate, typ.GameInstanceInit{})
}
