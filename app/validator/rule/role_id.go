package rule

import (
	"uwwolf/app/enum"
	"uwwolf/app/types"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slices"
)

const RoleIdTag = "role_id"

func ValidateRolePool(fl validator.FieldLevel) bool {
	var expectedRoleIds []types.RoleId

	if fl.Param() == "w" {
		expectedRoleIds = enum.WerewolfRoleIds
	} else {
		expectedRoleIds = enum.NonWerewolfRoleIds
	}

	if roleId, ok := fl.Field().Interface().(types.RoleId); ok {
		return slices.Contains(expectedRoleIds, roleId)
	} else if roleIds, ok := fl.Field().Interface().([]types.RoleId); ok {
		for _, roleId := range roleIds {
			if !slices.Contains(expectedRoleIds, roleId) {
				return false
			}
		}
	} else {
		return false
	}

	return true
}

func AddRolePoolTag(ut ut.Translator) error {
	message := "{0} is not valid"

	return ut.Add(RoleIdTag, message, true)
}

func RegisterRolePoolMessage(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T(RoleIdTag, fe.Field())

	return t
}
