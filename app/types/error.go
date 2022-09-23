package types

type ErrorTag = string

const (
	InvalidInputErrorTag ErrorTag = "INVALID_INPUT"
	UnauthorizedErrorTag ErrorTag = "UNAUTHORIZED"
	SystemErrorTag       ErrorTag = "SYSTEM_ERROR"
)
