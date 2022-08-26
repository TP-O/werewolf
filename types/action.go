package types

import "github.com/go-playground/validator/v10"

type ActionRequest struct {
	GameId    GameId
	Actor     PlayerId
	Targets   []PlayerId
	IsSkipped bool
}

type ActionResponse struct {
	Ok    bool
	Data  any
	Error *ErrorDetail
}

type ErrorDetail struct {
	Tag   ErrorTag
	Msg   validator.ValidationErrorsTranslations
	Alert string
}
