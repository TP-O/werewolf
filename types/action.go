package types

import "github.com/go-playground/validator/v10"

type ActionRequest struct {
	GameId    GameId     `validate:"required,number,gt=0"`
	ActorId   PlayerId   `validate:"required,len=28"`
	TargetIds []PlayerId `validate:"required,min=1,unique,dive,len=28"`
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
