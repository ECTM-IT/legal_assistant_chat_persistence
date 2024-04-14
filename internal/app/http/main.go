package http

import (
	"os"
	"runtime/debug"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/utils/env"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.uber.org/zap"
)

func Main() {
	logger := logs.Init()
	cfg := loadConfig()

	app, err := NewApplication(cfg, logger)
	if err != nil {
		trace := string(debug.Stack())
		app.logger.Warn("Failed to create application", zap.String("error", err.Error()), zap.String("trace", trace))
		os.Exit(1)
	}

	err = app.Run()
	if err != nil {
		trace := string(debug.Stack())
		app.logger.Warn("Application failed", zap.String("error", err.Error()), zap.String("trace", trace))
		os.Exit(1)
	}
}

func loadConfig() Config {
	var cfg Config
	cfg.BaseURL = env.GetString("BASE_URL", "http://0.0.0.0:4444")
	cfg.HTTPPort = env.GetInt("HTTP_PORT", 4444)
	cfg.Cookie.SecretKey = env.GetString("COOKIE_SECRET_KEY", "3iepwbkq5chsrusjoha26mnsjt233ujq")
	cfg.MongoDB.URI = env.GetString("MONGODB_URI", "mongodb+srv://alessiopersichetti:r9BtY7WjGv6ck5OS@latest.smobjvj.mongodb.net/?retryWrites=true&w=majority&appName=latest")
	cfg.MongoDB.Database = env.GetString("MONGODB_DATABASE", "legal_assistant")
	return cfg
}
