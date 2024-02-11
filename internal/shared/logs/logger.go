package logs

import (
	"os"

	"go.uber.org/zap"
)

// Logger defines the minimal logging interface
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Error(msg string, err error, fields ...zap.Field)
	// Extend with other levels as needed
}

// ZapLogger implements the Logger interface using Zap
type ZapLogger struct {
	*zap.Logger
}

func (l *ZapLogger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

func (l *ZapLogger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

func (l *ZapLogger) Error(msg string, err error, fields ...zap.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	l.Logger.Error(msg, fields...)
}

// Initialize the logger
func Init() Logger {
	logger, err := newZapLogger()
	if err != nil {
		panic(err) // Consider a more graceful error handling strategy
	}
	return &ZapLogger{logger}
}

func newZapLogger() (*zap.Logger, error) {
	env := os.Getenv("APP_ENV") // 'development' or 'production'
	var config zap.Config
	if env == "development" {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err // Properly handle error
	}
	return logger, nil
}
