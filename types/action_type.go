package types

import "github.com/go-playground/validator/v10"

type ActionData struct {
	GameId  string
	Actor   int
	Targets []int
	Skipped bool
	Payload []byte
}

type PerformResult struct {
	Ok       bool
	ErrorTag ErrorTag
	Errors   validator.ValidationErrorsTranslations
}
