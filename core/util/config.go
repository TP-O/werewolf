package util

import (
	"log"

	"github.com/spf13/viper"
)

type appConfig struct {
	Debug bool   `mapstructure:"debug"`
	Env   string `mapstructure:"env"`
	Port  uint16 `mapstructure:"port"`
}

type dbConfig struct {
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
	Keyspace string   `mapstructure:"keyspace"`
	Hosts    []string `mapstructure:"hosts"`
}

type gameConfig struct {
	MinCapacity           uint8  `mapstructure:"min_capacity"`
	MaxCapacity           uint8  `mapstructure:"max_capacity"`
	PreparationDuration   uint16 `mapstructure:"preparation_duration"`
	MinTurnDuration       uint16 `mapstructure:"min_turn_duration"`
	MaxTurnDuration       uint16 `mapstructure:"max_turn_duration"`
	MinDiscussionDuration uint16 `mapstructure:"min_discussion_duration"`
	MaxDiscussionDuration uint16 `mapstructure:"max_discussion_duration"`
}

type config struct {
	App  appConfig  `mapstructure:"app"`
	DB   dbConfig   `mapstructure:"database"`
	Game gameConfig `mapstructure:"game"`
}

var cfg *config

// Config returns config instance. Panic if config is not loaded.
func Config() *config {
	if cfg == nil {
		log.Panic("Load config first!")
	}
	return cfg
}

// loadDefaultConfig loads the default config values.
func loadDefaultConfig() {
	viper.SetDefault("app", map[string]interface{}{
		"debug": true,
		"env":   "development",
		"port":  8080,
	})
	viper.SetDefault("database", map[string]interface{}{
		"username": "username",
		"password": "password",
		"keyspace": "default_ns",
		"hosts":    []string{"locahost"},
	})
	viper.SetDefault("game", map[string]interface{}{
		"min_capacity":            5,
		"max_capacicty":           20,
		"preparation_time":        10,
		"min_turn_duration":       20,
		"max_turn_duration":       60,
		"min_discussion_duration": 60,
		"max_discussion_duration": 360,
	})
}

// LoadConfig loads config values from the given path and
// uses the default values if not set.
func LoadConfig(path string) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Unable to load config:", err)
		log.Println("Use default config!")
	}

	cfg = &config{}
	loadDefaultConfig()
	viper.Unmarshal(cfg)
}
