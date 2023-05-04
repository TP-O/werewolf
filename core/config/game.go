package config

import (
	"time"

	"github.com/spf13/viper"
)

type Game struct {
	MinCapacity           uint8         `mapstructure:"min_capacity"`
	MaxCapacity           uint8         `mapstructure:"max_capacity"`
	PreparationDuration   time.Duration `mapstructure:"preparation_duration"`
	MinTurnDuration       time.Duration `mapstructure:"min_turn_duration"`
	MaxTurnDuration       time.Duration `mapstructure:"max_turn_duration"`
	MinDiscussionDuration time.Duration `mapstructure:"min_discussion_duration"`
	MaxDiscussionDuration time.Duration `mapstructure:"max_discussion_duration"`
}

var _ configLoader = (*Game)(nil)

func (Game) loadDefault() {
	viper.SetDefault("game", map[string]interface{}{
		"min_capacity":            5,
		"max_capacicty":           20,
		"preparation_time":        "10s",
		"min_turn_duration":       "20s",
		"max_turn_duration":       "60s",
		"min_discussion_duration": "40s",
		"max_discussion_duration": "360s",
	})
}
