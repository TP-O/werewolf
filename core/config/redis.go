package config

import "github.com/spf13/viper"

type redis struct {
	Username          string   `mapstructure:"username"`
	Password          string   `mapstructure:"password"`
	MasterName        string   `mapstructure:"master_name"`
	SentinelAddresses []string `mapstructure:"sentinel_addresses"`
	PollSize          int      `mapstructure:"pool_size"`
}

var _ configLoader = (*redis)(nil)

func Redis() redis {
	return cfg.Redis
}

func (redis) loadDefault() {
	viper.SetDefault("redis", map[string]interface{}{
		"username":    "username",
		"password":    "password",
		"master_name": "mymaster",
		"poll_size":   10,
	})
}
