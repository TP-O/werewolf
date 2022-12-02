package validator

import (
	"reflect"
	"strings"
	"uwwolf/validator/rule"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
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

	// Rewrite struct fields as JSON
	validate.RegisterTagNameFunc(func(fl reflect.StructField) string {
		name := strings.SplitN(fl.Tag.Get("json"), ",", 2)

		return name[0]
	})

	// Customize rules for fileds
	validate.RegisterValidation(rule.CapacityTag, rule.ValidateCapacity)
	validate.RegisterValidation(rule.NumberOfWerewolvesTag, rule.ValidateNumberOfWerewolves)
	validate.RegisterValidation(rule.RoleIdTag, rule.ValidateRolePool)

	// Customize rules for structs
	// ...

	// Custom crror messages
	validate.RegisterTranslation(
		rule.CapacityTag,
		trans,
		rule.AddCapacityTag,
		rule.RegisterCapacityMessage,
	)
	validate.RegisterTranslation(
		rule.NumberOfWerewolvesTag,
		trans,
		rule.AddNumberOfWerewolvesTag,
		rule.RegisterNumberOfWerewolvesTagMessage,
	)
	validate.RegisterTranslation(
		rule.RoleIdTag,
		trans,
		rule.AddRolePoolTag,
		rule.RegisterRolePoolMessage,
	)
}

func SimpleValidateVar(data any, tag string) bool {
	return validate.Var(data, tag) != nil
}

func SimpleValidateStruct(data any) bool {
	return validate.Struct(data) != nil
}

func ValidateVar(data any, tag string) ValidationError {
	return handleError(validate.Var(data, tag))
}

func ValidateStruct(data any) ValidationError {
	return handleError(validate.Struct(data))
}

func handleError(err error) ValidationError {
	if err == nil {
		return nil
	}

	errs := err.(validator.ValidationErrors)

	return errs.Translate(trans)
}
