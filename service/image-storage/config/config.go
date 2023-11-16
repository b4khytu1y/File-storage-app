package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DB struct {
		Host     string
		User     string
		Password string
		DBName   string
		Port     string
	}
}

func LoadConfig(configPath string) (*Config, error) {
	cfg := Config{}

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.SetConfigFile("../../config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
