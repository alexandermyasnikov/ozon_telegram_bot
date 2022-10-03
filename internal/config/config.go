package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const ConfigFile = "./data/config.yaml"

type Config struct {
	Telegram struct {
		Token string
	}
}

func New(file string) (*Config, error) {
	var cfg Config

	rawYAML, err := os.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "parsing yaml")
	}

	return &cfg, nil
}

func (c *Config) TelegramToken() string {
	return c.Telegram.Token
}
