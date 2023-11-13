package config

import (
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
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
