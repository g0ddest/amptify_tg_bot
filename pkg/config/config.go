package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	TgToken string `yaml:"tg_token"`
}

func CreateConfig() (*Config, error) {
	config := &Config{}
	file, err := os.Open("./configs/config.yaml")
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
