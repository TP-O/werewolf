package types

import "github.com/go-playground/validator/v10"

type ActionRequest struct {
	GameId    GameId     `validate:"required,number,gt=0"`
	ActorId   PlayerId   `validate:"required,number,gt=0"`
	TargetIds []PlayerId `validate:"required,min=1,unique,dive,number,gt=0"`
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
