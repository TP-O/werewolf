package config

import (
	"uwwolf/util"

	"github.com/spf13/viper"
)

type databaseConfig struct {
	Host     string `mapstructure:"DB_HOST"`
	Port     uint   `mapstructure:"DB_PORT"`
	Username string `mapstructure:"DB_USERNAME"`
	Password string `mapstructure:"DB_PASSWORD"`
	Name     string `mapstructure:"DB_NAME"`
}

var Database *databaseConfig

func loadDatabaseConfig() {
	util.LoadDefaultConfigValues(map[string]interface{}{
		"DB_HOST":     "postgres",
		"DB_PORT":     5432,
		"DB_USERNAME": "username",
		"DB_PASSWORD": "password",
		"DB_NAME":     "db",
	})

	viper.Unmarshal(&Database)
}
