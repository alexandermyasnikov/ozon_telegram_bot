package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const ConfigFile = "./data/config.yaml"

type Config struct {
	Logger   LoggerConfig   `yaml:"logger"`
	Telegram TelegramConfig `yaml:"telegram"`
	Rates    RatesConfig    `yaml:"rates"`
	Database DatabaseConfig `yaml:"database"`
	Jaeger   JaegerConfig   `yaml:"jaeger"`
}

type LoggerConfig struct {
	Devel bool `yaml:"devel"`
}

type TelegramConfig struct {
	Enable bool   `yaml:"enable"`
	Token  string `yaml:"token"`
}

type RatesConfig struct {
	Service         string   `yaml:"service"`
	Base            string   `yaml:"base"`
	Codes           []string `yaml:"codes"`
	FreqUpdateInSec int      `yaml:"freqUpdateInSec"`
}

type DatabaseConfig struct {
	URL string `yaml:"url"`
}

type JaegerConfig struct {
	URL string `yaml:"url"`
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

func (c Config) GetLoggerDevel() bool {
	return c.Logger.Devel
}

func (c Config) TelegramEnable() bool {
	return c.Telegram.Enable
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
	return c.Rates.Service
}

func (c Config) GetDatabaseURL() string {
	return c.Database.URL
}

func (c Config) GetJaegerURL() string {
	return c.Jaeger.URL
}
