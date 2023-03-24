package config

import "github.com/spf13/viper"

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     uint16 `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Database string `mapstructure:"database"`
	PollSize int    `mapstructure:"pool_size"`
}

var _ configLoader = (*Postgres)(nil)

func (Postgres) loadDefault() {
	viper.SetDefault("postgres", map[string]interface{}{
		"host":      "postgres",
		"port":      5432,
		"username":  "username",
		"password":  "password",
		"database":  "db",
		"poll_size": 10,
	})
}
