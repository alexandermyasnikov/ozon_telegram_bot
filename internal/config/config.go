package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const ConfigFile = "./data/config.yaml"

type Config struct {
	Logger        LoggerConfig        `yaml:"logger"`
	Telegram      TelegramConfig      `yaml:"telegram"`
	Rates         RatesConfig         `yaml:"rates"`
	Database      DatabaseConfig      `yaml:"database"`
	Jaeger        JaegerConfig        `yaml:"jaeger"`
	CurrencyCache CacheConfig         `yaml:"currencyCache"`
	ReportCache   CacheConfig         `yaml:"reportCache"`
	Kafka         KafkaConfig         `yaml:"kafka"`
	Prometheus    PrometheusConfig    `yaml:"prometheus"`
	ReportService ReportServiceConfig `yaml:"reportService"`
}

type LoggerConfig struct {
	Devel bool `yaml:"devel"`
}

type TelegramConfig struct {
	Token string `yaml:"token"`
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

type CacheConfig struct {
	Enable bool `yaml:"enable"`
	Size   int  `yaml:"size"`
	TTL    int  `yaml:"ttl"`
}

type KafkaConfig struct {
	Addr string `yaml:"addr"`
}

type PrometheusConfig struct {
	Addr string `yaml:"addr"`
}

type ReportServiceConfig struct {
	Addr string `yaml:"addr"`
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

func (c Config) GetCurrencyCacheEnable() bool {
	return c.CurrencyCache.Enable
}

func (c Config) GetCurrencyCacheSize() int {
	return c.CurrencyCache.Size
}

func (c Config) GetCurrencyCacheTTL() int {
	return c.CurrencyCache.TTL
}

func (c Config) GetReportCacheEnable() bool {
	return c.ReportCache.Enable
}

func (c Config) GetReportCacheSize() int {
	return c.ReportCache.Size
}

func (c Config) GetReportCacheTTL() int {
	return c.ReportCache.TTL
}

func (c Config) GetKafkaAddr() string {
	return c.Kafka.Addr
}

func (c Config) GetPrometheusAddr() string {
	return c.Prometheus.Addr
}

func (c Config) GetReportServiceAddr() string {
	return c.ReportService.Addr
}
