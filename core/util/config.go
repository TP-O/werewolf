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
	Host     string `mapstructure:"host"`
	Port     uint16 `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	database string `mapstructure:"database"`
}

type redisConfig struct {
	Username          string   `mapstructure:"username"`
	Password          string   `mapstructure:"password"`
	MasterName        string   `mapstructure:"master_name"`
	SentinelAddresses []string `mapstructure:"sentinel_addresses"`
}

type gameConfig struct {
	MinCapacity           uint8         `mapstructure:"min_capacity"`
	MaxCapacity           uint8         `mapstructure:"max_capacity"`
	PreparationDuration   time.Duration `mapstructure:"preparation_duration"`
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
		"host":     "postgres",
		"port":     6379,
		"username": "username",
		"password": "password",
		"database": "default_ns",
	})
	viper.SetDefault("redis", map[string]interface{}{
		"username":    "username",
		"password":    "password",
		"master_name": "mymaster",
	})
	viper.SetDefault("game", map[string]interface{}{
		"min_capacity":            5,
		"max_capacicty":           20,
		"preparation_time":        "10s",
		"min_turn_duration":       "20s",
		"max_turn_duration":       "60s",
		"min_discussion_duration": "60s",
		"max_discussion_duration": "360s",
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
	if err := viper.Unmarshal(cfg); err != nil {
		panic(err)
	}

	viper.Unmarshal(cfg)
}
