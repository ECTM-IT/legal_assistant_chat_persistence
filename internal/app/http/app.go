package http

import (
	"os"
	"runtime/debug"

	"go.uber.org/zap"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/utils/env"
)

func Main() {
	// Initialize your ZapLogger
	logger := logs.Init()

	// Assuming Init() returns an instance of your ZapLogger which implements the Logger interface
	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Warn("Application failed", zap.String("error", err.Error()), zap.String("trace", trace))
		os.Exit(1)
	}
}

type Config struct {
	BaseURL  string
	HTTPPort int
	Cookie   struct {
		SecretKey string
	}
	MongoDB struct {
		URI      string
		Database string
	}
}

type Application struct {
	config Config
	logger logs.Logger
}

func run(logger logs.Logger) error {
	var cfg Config
	cfg.BaseURL = env.GetString("BASE_URL", "http://0.0.0.0:4444")
	cfg.HTTPPort = env.GetInt("HTTP_PORT", 4444)
	cfg.Cookie.SecretKey = env.GetString("COOKIE_SECRET_KEY", "3iepwbkq5chsrusjoha26mnsjt233ujq")
	cfg.MongoDB.URI = env.GetString("MONGODB_URI", "mongodb://0.0.0.0:27017/")
	cfg.MongoDB.Database = env.GetString("MONGODB_DATABASE", "legal_assistant")

	app := &Application{
		config: cfg,
		logger: logger,
	}

	return app.serveHTTP()
}
