package util

import "github.com/spf13/viper"

func LoadDefaultConfigValues[T any](configs map[string]T) {
	for key, value := range configs {
		LoadDefaultConfigValue(key, value)
	}
}

func LoadDefaultConfigValue[T any](key string, value T) {
	if viper.Get(key) == nil || viper.Get(key) == "" {
		viper.Set(key, value)
	}
}
