package http

import (
	"os"
	"runtime/debug"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/utils/env"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// Main is the entry point of the application.
func Main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		// Log warning but don't exit - we might be running in production where .env is not needed
		logger := logs.Init()
		logger.Warn("No .env file found", zap.Error(err))
		logger.Sync()
	}

	logger := logs.Init()
	defer logger.Sync()

	cfg := loadConfig(logger)
	app, err := NewApplication(cfg, logger)
	if err != nil {
		handleErrorAndExit(logger, "Failed to create application", err)
	}

	if err = app.Run(); err != nil {
		handleErrorAndExit(logger, "Application failed", err)
	}
}

// loadConfig loads the configuration settings from the environment.
func loadConfig(logger logs.Logger) Config {
	var cfg Config
	cfg.BaseURL = env.GetString("BASE_URL", "http://0.0.0.0:4444")
	cfg.HTTPPort = env.GetInt("HTTP_PORT", 4444)
	cfg.Cookie.SecretKey = env.GetString("COOKIE_SECRET_KEY", "3iepwbkq5chsrusjoha26mnsjt233ujq")
	cfg.MongoDB.URI = env.GetString("MONGODB_URI", "mongodb://legal-assist-chat-persistence-dev-qa:AFWmqIy38KoY8pquvfuZmNFmEdUoYJ139dQAYmNhS7urHjVjOVDliCgNMfQhCOED6m1DJ3uxi8k3ACDbvUuSog%3D%3D@legal-assist-chat-persistence-dev-qa.mongo.cosmos.azure.com:10255/?ssl=true&replicaSet=globaldb&retrywrites=false&maxIdleTimeMS=120000")
	cfg.MongoDB.Database = env.GetString("MONGODB_DATABASE", "legal_assistant")

	logger.Info("Configuration loaded",
		zap.String("baseURL", cfg.BaseURL),
		zap.Int("httpPort", cfg.HTTPPort),
		zap.String("mongoDBURI", cfg.MongoDB.URI),
		zap.String("mongoDBDatabase", cfg.MongoDB.Database),
	)

	return cfg
}

// handleErrorAndExit logs the error and exits the application.
func handleErrorAndExit(logger logs.Logger, message string, err error) {
	trace := string(debug.Stack())
	logger.Error(message, err, zap.String("trace", trace))
	os.Exit(1)
}
