package validation

import (
	"reflect"
	"strings"
	"uwwolf/app/dto"
	"uwwolf/app/validation/strct"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	trans ut.Translator
)

func ImproveValidator(validate *validator.Validate) {
	uni := ut.New(en.New())
	trans, _ = uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans) // nolint errcheck

	// Rewrite struct field to JSON tag
	validate.RegisterTagNameFunc(func(fl reflect.StructField) string {
		name := strings.SplitN(fl.Tag.Get("json"), ",", 2)
		return name[0]
	})

	addCutomizedValidationTags(trans, validate)
	addCustomizedFieldRules(validate)
	addCustomizedStructRules(validate)
}

func addCutomizedValidationTags(trans ut.Translator, validate *validator.Validate) {
	validate.RegisterTranslation( // nolint errcheck
		strct.TurnDurationTag,
		trans,
		strct.AddTurnDurationTag,
		strct.RegisterTurnDurationMessage,
	)
	validate.RegisterTranslation( // nolint errcheck
		strct.DiscussionDurationTag,
		trans,
		strct.AddDiscussionDurationTag,
		strct.RegisterDiscussionDurationMessage,
	)
	validate.RegisterTranslation( // nolint errcheck
		strct.RoleIDTag,
		trans,
		strct.AddRoleIDsTag,
		strct.RegisterRoleIDsMessage,
	)
	validate.RegisterTranslation( // nolint errcheck
		strct.RequiredRoleIDsTag,
		trans,
		strct.AddRequiredRoleIDsTag,
		strct.RegisterRequiredRoleIDsMessage,
	)
}

func addCustomizedFieldRules(validate *validator.Validate) {
	//
}

func addCustomizedStructRules(validate *validator.Validate) {
	validate.RegisterStructValidation(strct.ValidateGameConfig, dto.ReplaceGameConfigDto{})
}

type validationErrors = map[string]string

func FormatValidationError(err error) validationErrors {
	fieldErrors := make(validationErrors)

	errs, ok := err.(validator.ValidationErrors)
	if errs == nil || !ok {
		fieldErrors["*"] = "Please double check the json format!"
		return fieldErrors
	}

	for _, e := range errs {
		fieldErrors[e.Field()] = e.Translate(trans)
	}

	return fieldErrors
}
