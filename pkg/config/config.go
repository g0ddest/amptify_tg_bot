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
	if err == nil {
		defer func(file *os.File) {
			_ = file.Close()
		}(file)

		d := yaml.NewDecoder(file)
		if err := d.Decode(&config); err != nil {
			return nil, err
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	if config.TgToken == "" {
		config.TgToken = os.Getenv("TG_TOKEN")
	}

	if config.TgToken == "" {
		return nil, ErrMissingTgToken
	}

	return config, nil
}

var ErrMissingTgToken = &MissingTgTokenError{}

type MissingTgTokenError struct{}

func (e *MissingTgTokenError) Error() string {
	return "tg_token is missing: not found in config file or environment variable TG_TOKEN"
}
