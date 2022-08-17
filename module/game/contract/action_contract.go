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

	// Set action's state. Return false if type conversion is failed.
	SetState(state any) bool

	// Validate action's input first, then execute it if the
	// validation is successful. Only supposed to fail if
	// and only if an error message is returned.
	Perform(data *types.ActionData) *types.PerformResult

	// Validate the action's input. Each action has different rules
	// for data validation.
	Validate(data *types.ActionData) validator.ValidationErrorsTranslations

	// Execute the action with receied data. Return the result of execution
	// and error message, if any. The execution is only supposed to fail if
	// and only if an error message is returned. The first response arg is
	// just a status of the execution, so its meaning depends on contenxt.
	Execute(data *types.ActionData) *types.PerformResult
}
