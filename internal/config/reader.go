package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.AllowEmptyEnv(true)

	if err := viper.ReadInConfig(); err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return nil, fmt.Errorf("reading config: %w", err)
		}
	}
	viper.AutomaticEnv()

	return &Config{
		AppEnv:   AppEnv(loadString("APP_ENV")),
		TimeZone: loadString("TZ"),
		AppDebug: loadBool("APP_DEBUG"),
		Mysql: Mysql{
			Host:     loadString("MYSQL_HOST"),
			Port:     loadString("MYSQL_PORT"),
			UserName: loadString("MYSQL_USER"),
			Password: loadString("MYSQL_PASSWORD"),
			Database: loadString("MYSQL_DB"),
		},
		HTTP: HTTP{
			Host: loadString("APP_HTTP_HOST"),
			Port: loadInt("APP_HTTP_PORT"),
		},
		Redis: Redis{
			Host:     loadString("REDIS_HOST"),
			Port:     loadString("REDIS_PORT"),
			Database: loadInt("REDIS_DATABASE"),
			Password: loadString("REDIS_PASSWORD"),
		},
	}, nil
}
