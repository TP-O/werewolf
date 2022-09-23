package rule

import (
	"uwwolf/app/model"
	"uwwolf/app/types"
	"uwwolf/database"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

const RolePoolTag = "role_pool"

func ValidateRolePool(fl validator.FieldLevel) bool {
	var roles []*model.Role
	roleIds := fl.Field().Interface().([]types.RoleId)
	database.Client().Where("id IN (?)", roleIds).Find(&roles)

	if len(roles) != len(roleIds) {
		return false
	}

	return true
}

func AddRolePoolTag(ut ut.Translator) error {
	message := "{0} is not valid"

	return ut.Add(RolePoolTag, message, true)
}

func RegisterRolePoolMessage(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(RolePoolTag, fe.Field())

	return t
}
