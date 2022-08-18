package contract

import (
	"github.com/go-playground/validator/v10"

	"uwwolf/types"
)

type Action interface {
	// Get action's name.
	GetName() string

	// Export action's state  as JSON string.
	JsonState() string

	// Validate action's input first, then execute it if the
	// validation is successful. Only supposed to fail if
	// and only if an error message is returned.
	Perform(req *types.ActionRequest) *types.ActionResponse

	// Validate the action's input. Each action has different rules
	// for data validation.
	Validate(req *types.ActionRequest) validator.ValidationErrorsTranslations

	// Execute the action with receied data. Return the result of execution
	// and error message, if any. The execution is only supposed to fail if
	// and only if an error message is returned. The first response arg is
	// just a status of the execution, so its meaning depends on contenxt.
	Execute(req *types.ActionRequest) *types.ActionResponse
}
