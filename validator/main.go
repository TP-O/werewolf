package validator

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"

	"uwwolf/module/game/types"
	"uwwolf/validator/rule"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	en := en.New()
	uni = ut.New(en, en)
	trans, _ = uni.GetTranslator("en")

	validate = validator.New()

	en_translations.RegisterDefaultTranslations(validate, trans)

	// Customize rules for fileds
	validate.RegisterValidation(rule.GameCapacityTag, rule.GameCapacityValidate)

	// Customize rules for structs
	validate.RegisterStructValidation(rule.GameInstanceInitValidate, types.GameInstanceInit{})

	// Custom crror messages
	validate.RegisterTranslation(
		rule.GameCapacityTag,
		trans,
		rule.AddGameCapacityTag,
		rule.RegisterGameCapacityMessage,
	)
	validate.RegisterTranslation(
		rule.NumberOfWerewolvesTag,
		trans,
		func(ut ut.Translator) error { return rule.AddNumberOfWerewolvesTag(ut) },
		rule.RegisterNumberOfWerewolvesTagMessage,
	)
}

func SimpleValidateVar(data interface{}, tag string) bool {
	return validate.Var(data, tag) != nil
}

func SimpleValidateStruct(data interface{}, tag ...string) bool {
	return validate.Struct(data) != nil
}

func ValidateVar(data interface{}, tag string) validator.ValidationErrorsTranslations {
	return handleError(validate.Var(data, tag))
}

func ValidateStruct(data interface{}, tag ...string) validator.ValidationErrorsTranslations {
	return handleError(validate.Struct(data))
}

func handleError(err error) validator.ValidationErrorsTranslations {
	if err == nil {
		return nil
	}

	errs := err.(validator.ValidationErrors)

	return errs.Translate(trans)
}
