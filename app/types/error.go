package types

import "github.com/go-playground/validator/v10"

type ErrorTag = string

type ErrorDetail struct {
	Tag   ErrorTag
	Msg   validator.ValidationErrorsTranslations
	Alert string
}
