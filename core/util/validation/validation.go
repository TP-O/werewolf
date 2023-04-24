package validation

import (
	"reflect"
	"strings"
	"uwwolf/config"
	"uwwolf/server/dto"
	"uwwolf/util/validation/strct"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	trans ut.Translator
)

func Setup(validate *validator.Validate) {
	uni := ut.New(en.New())
	trans, _ = uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, trans) // nolint errcheck

	// Rewrite struct field to JSON tag
	validate.RegisterTagNameFunc(func(fl reflect.StructField) string {
		name := strings.SplitN(fl.Tag.Get("json"), ",", 2)
		return name[0]
	})

	config := config.Load("../..").Game
	addCutomizedValidationTags(trans, validate, config)
	addCustomizedFieldRules(validate)
	addCustomizedStructRules(validate, config)
}

func addCutomizedValidationTags(trans ut.Translator, validate *validator.Validate, config config.Game) {
	validate.RegisterTranslation( // nolint errcheck
		strct.TurnDurationTag,
		trans,
		strct.AddTurnDurationTag(config),
		strct.RegisterTurnDurationMessage,
	)
	validate.RegisterTranslation( // nolint errcheck
		strct.DiscussionDurationTag,
		trans,
		strct.AddDiscussionDurationTag(config),
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

func addCustomizedStructRules(validate *validator.Validate, config config.Game) {
	validate.RegisterStructValidation(strct.ValidateGameConfig(config), dto.ReplaceGameConfigDto{})
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
