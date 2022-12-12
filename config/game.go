package config

import (
	"github.com/spf13/viper"
)

type gameConfig struct {
	MinCapacity           uint8  `mapstructure:"min_capacity"`
	MaxCapacity           uint8  `mapstructure:"max_capacity"`
	MinPollCapacity       uint8  `mapstructure:"min_poll_capacity"`
	PreparationTime       uint16 `mapstructure:"preparation_time"`
	MinTurnDuration       uint16 `mapstructure:"min_turn_duration"`
	MaxTurnDuration       uint16 `mapstructure:"max_turn_duration"`
	MinDiscussionDuration uint16 `mapstructure:"min_discussion_duration"`
	MaxDiscussionDuration uint16 `mapstructure:"max_discussion_duration"`
}

func loadDefaultGame() {
	viper.SetDefault("game", map[string]interface{}{
		"min_capacity":            5,
		"max_capacicty":           20,
		"min_poll_capacity":       1,
		"preparation_time":        10,
		"min_turn_duration":       20,
		"max_turn_duration":       60,
		"min_discussion_duration": 60,
		"max_discussion_duration": 360,
	})
}
