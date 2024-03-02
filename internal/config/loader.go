package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func loadString(envName string) string {
	validate(envName)
	return viper.GetString(envName)
}

func loadInt(envName string) int {
	validate(envName)
	return viper.GetInt(envName)
}

func loadInt64(envName string) int64 {
	validate(envName)
	return viper.GetInt64(envName)
}

func loadBool(envName string) bool {
	validate(envName)
	return viper.GetBool(envName)
}

func loadFloat64(envName string) float64 {
	validate(envName)
	return viper.GetFloat64(envName)
}

func validate(envName string) {
	isSet := viper.IsSet(envName)
	if !isSet {
		panic(fmt.Sprintf("environment variable [%s] does not exist", envName))
	}
}
