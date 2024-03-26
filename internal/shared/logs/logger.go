package logs

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, err error, fields ...zap.Field)
}

type ZapLogger struct {
	*zap.Logger
}

func (l *ZapLogger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

func (l *ZapLogger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

func (l *ZapLogger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

func (l *ZapLogger) Error(msg string, err error, fields ...zap.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	l.Logger.Error(msg, fields...)
}

// Init initializes the logger
func Init() Logger {
	logger, err := newZapLogger()
	if err != nil {
		panic(err)
	}
	return &ZapLogger{logger}
}

func newZapLogger() (*zap.Logger, error) {
	env := os.Getenv("APP_ENV") // 'development' or 'production'

	var config zap.Config
	if env == "development" {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
	}

	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if env == "development" {
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}
