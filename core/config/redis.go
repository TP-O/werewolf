package config

import "github.com/spf13/viper"

type Redis struct {
	Username          string   `mapstructure:"username"`
	Password          string   `mapstructure:"password"`
	MasterName        string   `mapstructure:"master_name"`
	SentinelAddresses []string `mapstructure:"sentinel_addresses"`
	PollSize          int      `mapstructure:"pool_size"`
}

var _ configLoader = (*Redis)(nil)

func (Redis) loadDefault() {
	viper.SetDefault("redis", map[string]interface{}{
		"username":    "username",
		"password":    "password",
		"master_name": "mymaster",
		"poll_size":   10,
	})
}
