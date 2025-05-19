package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Dev        bool                    `yaml:"dev"`
	Encryption EncryptionConfiguration `yaml:"encryption"`
	Database   DatabaseConfiguration   `yaml:"database"`
}

func GenerateConfigFromFile(configurationFile string) (*Config, error) {
	var conf Config

	openedFile, err := os.ReadFile(configurationFile)
	if err != nil {
		return &conf, err
	}

	if err := yaml.Unmarshal(openedFile, &conf); err != nil {
		return &conf, err
	}

	validate := validator.New()
	err = validate.Struct(conf)
	if err != nil {
		return &conf, err
	}

	return &conf, nil
}
