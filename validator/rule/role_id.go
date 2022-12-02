package rule

import (
	"uwwolf/game/enum"
	"uwwolf/game/types"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slices"
)

const RoleIdTag = "role_id"

func ValidateRolePool(fl validator.FieldLevel) bool {
	var expectedRoleIds []enum.RoleID

	if fl.Param() == "w" {
		expectedRoleIds = types.RoleIDsByFactionID[enum.WerewolfFactionID]
	} else {
		expectedRoleIds = types.RoleIDsByFactionID[enum.VillagerFactionID]
	}

	if roleId, ok := fl.Field().Interface().(enum.RoleID); ok {
		return slices.Contains(expectedRoleIds, roleId)
	} else if roleIds, ok := fl.Field().Interface().([]enum.RoleID); ok {
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
