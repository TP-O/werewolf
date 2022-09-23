package config

import (
	"github.com/spf13/viper"

	"uwwolf/app/util"
)

type cacheConfig struct {
	LocalTTL string `mapstructure:"LOCAL_CACHE_TTL"`
}

func (c *cacheConfig) load() {
	util.LoadDefaultConfigValues(map[string]any{
		"LOCAL_CACHE_TTL": "60m",
	})

	viper.Unmarshal(&c)
}
