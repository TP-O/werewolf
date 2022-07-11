package rule

import (
	"github.com/go-playground/validator/v10"

	"uwwolf/config"
)

const GameCapacityRule = "game_capacity"

func GameCapacityValidate(fl validator.FieldLevel) bool {
	return fl.Field().Uint() <= uint64(config.Game.MaxCapacity) &&
		fl.Field().Uint() >= uint64(config.Game.MinCapacity)
}
