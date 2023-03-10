package util

import (
	"log"
	"time"

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

type redisConfig struct {
	Hosts    []string `mapstructure:"hosts"`
	Password string   `mapstructure:"password"`
}

type gameConfig struct {
	MinCapacity           uint8         `mapstructure:"min_capacity"`
	MaxCapacity           uint8         `mapstructure:"max_capacity"`
	PreparationDuration   uint16        `mapstructure:"preparation_duration"`
	MinTurnDuration       time.Duration `mapstructure:"min_turn_duration"`
	MaxTurnDuration       time.Duration `mapstructure:"max_turn_duration"`
	MinDiscussionDuration time.Duration `mapstructure:"min_discussion_duration"`
	MaxDiscussionDuration time.Duration `mapstructure:"max_discussion_duration"`
}

type config struct {
	App   appConfig   `mapstructure:"app"`
	DB    dbConfig    `mapstructure:"database"`
	Game  gameConfig  `mapstructure:"game"`
	Redis redisConfig `mapstructure:"redis"`
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
		"hosts":    []string{"locahost:6379"},
	})
	viper.SetDefault("game", map[string]interface{}{
		"min_capacity":            5,
		"max_capacicty":           20,
		"preparation_time":        10,
		"min_turn_duration":       "20s",
		"max_turn_duration":       "60s",
		"min_discussion_duration": "60s",
		"max_discussion_duration": "360s",
	})

	viper.SetDefault("redis", map[string]interface{}{
		"password": "password",
		"hosts":    []string{"locahost:6379"},
	})
}

// LoadConfig loads config values from the given path and
// uses the default values if not set.
func LoadConfig(path string) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Unable to load config:", err)
		log.Println("Use default config!")
	}

	cfg = &config{}
	loadDefaultConfig()
	viper.Unmarshal(cfg)
}
