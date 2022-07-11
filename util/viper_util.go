package util

import "github.com/spf13/viper"

func LoadDefaultConfigValues(configs map[string]interface{}) {
	for key, value := range configs {
		LoadDefaultConfigValue(key, value)
	}
}

func LoadDefaultConfigValue(key string, value interface{}) {
	if viper.Get(key) == "" {
		viper.Set(key, value)
	}
}
