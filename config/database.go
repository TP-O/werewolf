package config

import (
	"github.com/spf13/viper"
)

type dbConfig struct {
	Username string   `mapstructure:"username"`
	Password string   `mapstructure:"password"`
	Keyspace string   `mapstructure:"keyspace"`
	Hosts    []string `mapstructure:"hosts"`
}

func loadDefaultDB() {
	viper.SetDefault("database", map[string]interface{}{
		"username": "username",
		"password": "password",
		"keyspace": "default_ns",
		"hosts":    []string{"locahost"},
	})
}
