package tag

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

const RoleIDTag = "role_id"

func AddRoleIDsTag(ut ut.Translator) error {
	message := "{0} contains invalid role id"

	return ut.Add(RoleIDTag, message, true)
}

func RegisterRoleIDsMessage(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(RoleIDTag, fe.Field())

	return t
}
