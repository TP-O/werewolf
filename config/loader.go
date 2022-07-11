package config

import (
	"github.com/spf13/viper"
)

func LoadConfigs() error {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	loadAppConfig()
	loadDatabaseConfig()
	loadGameConfig()

	return nil
}
