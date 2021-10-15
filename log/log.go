package log

import (
	"github.com/tyuhara/team-review/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a new zap logger with the given log level, service name and environment.
func New(conf *config.Config) (*zap.Logger, error) {
	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(conf.LogLevel),
		Development:       false,
		Encoding:          "json",
		DisableCaller:     true,
		DisableStacktrace: true,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			MessageKey:     "message",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.EpochMillisTimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	return zapConfig.Build()
}

// NewDiscard creates logger which output to ioutil.Discard.
// This can be used for testing.
func NewDiscard() *zap.Logger {
	return zap.NewNop()
}
