package validator

import (
	"reflect"
	"strings"
	"uwwolf/validator/tag"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	uni = ut.New(en.New())
	trans, _ = uni.GetTranslator("en")
	validate = validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)

	// Rewrite struct field as JSON tag
	validate.RegisterTagNameFunc(func(fl reflect.StructField) string {
		name := strings.SplitN(fl.Tag.Get("json"), ",", 2)

		return name[0]
	})

	addCutomizedValidationTags(validate)
	addCustomizedFieldRules(validate)
	addCustomizedStructRules(validate)
}

func addCutomizedValidationTags(validate *validator.Validate) {
	validate.RegisterTranslation(
		tag.TurnDurationTag,
		trans,
		tag.AddTurnDurationTag,
		tag.RegisterTurnDurationMessage,
	)
	validate.RegisterTranslation(
		tag.DiscussionDurationTag,
		trans,
		tag.AddDiscussionDurationTag,
		tag.RegisterDiscussionDurationMessage,
	)
	validate.RegisterTranslation(
		tag.GameCapacityTag,
		trans,
		tag.AddGameCapacityTag,
		tag.RegisterGameCapacityMessage,
	)
	validate.RegisterTranslation(
		tag.RoleIDTag,
		trans,
		tag.AddRoleIDsTag,
		tag.RegisterRoleIDsMessage,
	)
	validate.RegisterTranslation(
		tag.NumberWerewolvesTag,
		trans,
		tag.AddNumberWerewolvesTag,
		tag.RegisterNumberWerewolvesMessage,
	)
}

func addCustomizedFieldRules(validate *validator.Validate) {
	//
}

func addCustomizedStructRules(validate *validator.Validate) {
	// validate.RegisterStructValidation(strct.ValidateGameSetting, types.GameSetting{})
}

func SimpleValidateVar(data any, tag string) bool {
	return validate.Var(data, tag) != nil
}

func SimpleValidateStruct(data any) bool {
	return validate.Struct(data) != nil
}

func ValidateVar(data any, tag string) *ValidationErrors {
	return handleError(validate.Var(data, tag))
}

func ValidateStruct(data any) *ValidationErrors {
	return handleError(validate.Struct(data))
}

func handleError(err error) *ValidationErrors {
	fieldErrors := new(ValidationErrors)

	if err == nil {
		return fieldErrors
	}

	for _, e := range err.(validator.ValidationErrors) {
		fieldErrors.FieldViolations = append(
			fieldErrors.FieldViolations,
			&errdetails.BadRequest_FieldViolation{
				Field:       e.Field(),
				Description: e.Translate(trans),
			})
	}

	return fieldErrors
}
