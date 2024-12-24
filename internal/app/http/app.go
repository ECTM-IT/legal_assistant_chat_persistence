package http

import (
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/db"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/mongo"
)

// Config holds the configuration settings for the application.
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

// Application holds the configuration, logger, and services for the application.
type Application struct {
	config   Config
	logger   logs.Logger
	services *db.Services
	client   *mongo.Client
}

// NewApplication initializes a new Application with the provided configuration and logger.
func NewApplication(cfg Config, logger logs.Logger) (*Application, error) {
	client, err := db.Connect(cfg.MongoDB.URI, 60*time.Second, logger)
	if err != nil {
		return nil, err
	}

	laDatabase := db.CreateDB(client, cfg.MongoDB.Database, logger)
	services := db.InitializeServices(laDatabase, logger)

	return &Application{
		config:   cfg,
		logger:   logger,
		services: services,
		client:   client,
	}, nil
}

// Run starts the HTTP server for the application.
func (app *Application) Run() error {
	server := NewHTTPServer(&app.config, app.services, app.logger)
	return server.ServeHTTP()
}
