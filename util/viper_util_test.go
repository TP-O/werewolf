package util_test

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"uwwolf/util"
)

const (
	env1 = "environment 01"
	env2 = "environment 02"
)

func TestLoadDefaultConfigValues(t *testing.T) {
	//=============================================================
	// Load env variables to viper
	util.LoadDefaultConfigValues(map[string]string{
		"ENV_1": env1,
		"ENV_2": env2,
	})

	assert.Equal(t, viper.Get("ENV_1"), env1)
	assert.Equal(t, viper.Get("ENV_2"), env2)

	//=============================================================
	// Map to struct
	type config struct {
		Env1 string `mapstructure:"ENV_1"`
		Env2 string `mapstructure:"ENV_2"`
	}

	var cfg config
	viper.Unmarshal(&cfg)

	assert.Equal(t, cfg.Env1, env1)
	assert.Equal(t, cfg.Env2, env2)
}
func TestLoadDefaultConfigValue(t *testing.T) {
	//=============================================================
	// Load env variable to viper
	util.LoadDefaultConfigValue("ENV_1", env1)

	assert.Equal(t, viper.Get("ENV_1"), env1)

	//=============================================================
	// Map to struct
	type config struct {
		Env1 string `mapstructure:"ENV_1"`
	}

	var cfg config
	viper.Unmarshal(&cfg)

	assert.Equal(t, cfg.Env1, env1)
}
