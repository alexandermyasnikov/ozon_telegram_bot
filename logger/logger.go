package logger

import (
	"log"

	"go.uber.org/zap"
)

var _logger *zap.Logger //nolint:gochecknoglobals

func init() { //nolint:gochecknoinits
	InitLogger(false)
}

func InitLogger(develMode bool) {
	var err error

	if develMode {
		_logger, err = zap.NewDevelopment()
	} else {
		cfg := zap.NewProductionConfig()
		cfg.DisableCaller = true
		cfg.DisableStacktrace = true
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		_logger, err = cfg.Build()
	}

	if err != nil {
		log.Fatal("cannot init zap", err)
	}
}

func Infof(template string, args ...interface{}) {
	_logger.Sugar().Infof(template, args...)
}

func Errorf(template string, args ...interface{}) {
	_logger.Sugar().Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	_logger.Sugar().Fatalf(template, args...)
}
