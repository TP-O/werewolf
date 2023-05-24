package config

import (
	"time"

	"github.com/spf13/viper"
)

type Game struct {
	MinCapacity           uint8         `mapstructure:"minCapacity"`
	MaxCapacity           uint8         `mapstructure:"maxCapacity"`
	PreparationDuration   time.Duration `mapstructure:"preparationDuration"`
	MinTurnDuration       time.Duration `mapstructure:"minTurnDuration"`
	MaxTurnDuration       time.Duration `mapstructure:"maxTurnDuration"`
	MinDiscussionDuration time.Duration `mapstructure:"minDiscussionDuration"`
	MaxDiscussionDuration time.Duration `mapstructure:"maxDiscussionDuration"`
}

var _ configLoader = (*Game)(nil)

func (Game) loadDefault() {
	viper.SetDefault("game", map[string]interface{}{
		"minCapacity":           5,
		"maxCapacicty":          20,
		"preparationxTime":      "10s",
		"minTurnDuration":       "20s",
		"maxTurnDuration":       "60s",
		"minDiscussionDuration": "40s",
		"maxDiscussionDuration": "360s",
	})
}
