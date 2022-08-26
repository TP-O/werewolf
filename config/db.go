package config

import (
	"uwwolf/util"

	"github.com/spf13/viper"
)

type dbConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     int    `mapstructure:"DB_PORT"`
	Username string `mapstructure:"DB_USERNAME"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
}

func (c *dbConfig) load() {
	util.LoadDefaultConfigValues(map[string]any{
		"DB_HOST":     "postgres",
		"DB_PORT":     5432,
		"DB_USERNAME": "username",
		"DB_PASSWORD": "password",
		"DB_NAME":     "db",
	})

	viper.Unmarshal(&c)
}
