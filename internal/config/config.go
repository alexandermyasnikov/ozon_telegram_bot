package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const ConfigFile = "./data/config.yaml"

type Config struct {
	Telegram TelegramConfig `yaml:"telegram"`
	Rates    RatesConfig    `yaml:"rates"`
}

type TelegramConfig struct {
	Token string `yaml:"token"`
}

type RatesConfig struct {
	service         string   `yaml:"service"`
	Base            string   `yaml:"base"`
	Codes           []string `yaml:"codes"`
	FreqUpdateInSec int      `yaml:"freqUpdateInSec"`
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

func (c Config) TelegramToken() string {
	return c.Telegram.Token
}
func (c Config) GetBaseCurrencyCode() string {
	return c.Rates.Base
}

func (c Config) GetCurrencyCodes() []string {
	return c.Rates.Codes
}

func (c Config) GetFrequencyRateUpdateSec() int {
	return c.Rates.FreqUpdateInSec
}

func (c Config) GetRatesService() string {
	return c.Rates.service
}
