package types

type ErrorTag = string

const (
	InvalidInputErrorTag ErrorTag = "INVALID_INPUT"
	UnauthorizedErrorTag ErrorTag = "UNAUTHORIZED"
	SystemErrorTag       ErrorTag = "SYSTEM_ERROR"
)

type ErrorProperty = string

const (
	SystemErrorProperty ErrorProperty = "system"
	AlertErrorPropery   ErrorProperty = "alert"
)
