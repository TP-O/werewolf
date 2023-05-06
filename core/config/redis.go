package config

import "github.com/spf13/viper"

type Redis struct {
	Username          string   `mapstructure:"username"`
	Password          string   `mapstructure:"password"`
	MasterName        string   `mapstructure:"masterName"`
	SentinelAddresses []string `mapstructure:"sentinelAddresses"`
	PollSize          int      `mapstructure:"poolSize"`
}

var _ configLoader = (*Redis)(nil)

func (Redis) loadDefault() {
	viper.SetDefault("redis", map[string]interface{}{
		"username":   "username",
		"password":   "password",
		"masterName": "mymaster",
		"pollSize":   10,
	})
}
