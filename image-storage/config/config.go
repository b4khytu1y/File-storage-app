package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DB struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		DBHost   string `yaml:"host"`
		DBPort   string `yaml:"port"`
	} `yaml:"db"`
}

func LoadConfig(configPath string) (*Config, error) {
	cfg := &Config{}
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(bytes, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
