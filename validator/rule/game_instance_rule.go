package rule

import (
	"math"
	"uwwolf/database"
	"uwwolf/game/contract"
	"uwwolf/model"

	"github.com/go-playground/validator/v10"
)

const (
	InvalidRoleIds            = "role_ids"
	InvalidNumberOfWerewolves = "number_of_werewolves"
)

func GameInstanceInitValidate(sl validator.StructLevel) {
	input := sl.Current().Interface().(contract.GameInstanceInit)

	roleIdsValidate(sl, input.RolePool)
	numberOfWerewolvesValidate(sl, input.Capacity, input.NumberOfWerewolves)
}

func roleIdsValidate(sl validator.StructLevel, rolePool []uint) {
	var counter int64
	database.DB().
		Model(&model.Role{}).
		Where("id IN (?)", rolePool).
		Count(&counter)

	if int(counter) != len(rolePool) {
		sl.ReportError(nil, "RolePool", "RolePool", InvalidRoleIds, "")
	}
}

func numberOfWerewolvesValidate(sl validator.StructLevel, capacity uint, numberOfWerewolves uint) {
	maxNumberOfWerewolves := uint(math.Round(float64(capacity)/2)) - 2

	if numberOfWerewolves > maxNumberOfWerewolves {
		sl.ReportError(nil, "NumberOfWerewolves", "NumberOfWerewolves", InvalidNumberOfWerewolves, "")

		return
	}
}
