package util

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const (
	env1 = "environment 01"
	env2 = "environment 02"
)

func TestLoadDefaultConfigValues(t *testing.T) {
	type config struct {
		Env1 string `mapstructure:"ENV_1"`
		Env2 string `mapstructure:"ENV_2"`
	}

	LoadDefaultConfigValues(map[string]string{
		"ENV_1": env1,
		"ENV_2": env2,
	})

	assert.Equal(t, viper.Get("ENV_1"), env1)
	assert.Equal(t, viper.Get("ENV_2"), env2)

	var cfg config
	viper.Unmarshal(&cfg)

	assert.Equal(t, cfg.Env1, env1)
	assert.Equal(t, cfg.Env2, env2)
}
func TestLoadDefaultConfigValue(t *testing.T) {
	type config struct {
		Env1 string `mapstructure:"ENV_1"`
	}

	LoadDefaultConfigValue("ENV_1", env1)

	assert.Equal(t, viper.Get("ENV_1"), env1)

	var cfg config
	viper.Unmarshal(&cfg)

	assert.Equal(t, cfg.Env1, env1)
}
