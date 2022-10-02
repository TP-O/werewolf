package types

import "github.com/go-playground/validator/v10"

type ActionId = int

type ActionRequest struct {
	GameId    GameId     `json:"gameId" validate:"required,number,gt=0"`
	ActorId   PlayerId   `json:"actorId" validate:"required,len=28"`
	TargetIds []PlayerId `json:"targetIds" validate:"required,min=1,unique,dive,len=28"`
	IsSkipped bool       `json:"isSkipped"`
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
