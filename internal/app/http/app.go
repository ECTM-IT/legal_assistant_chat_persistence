package http

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/db"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
)

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
	config   Config
	logger   logs.Logger
	services *db.Services
}

func NewApplication(cfg Config, logger logs.Logger) (*Application, error) {
	client, err := db.Connect(cfg.MongoDB.URI, 60)
	if err != nil {
		return nil, err
	}

	laDatabase := db.CreateDB(client)
	services := db.InitializeServices(laDatabase)

	return &Application{
		config:   cfg,
		logger:   logger,
		services: services,
	}, nil
}

func (app *Application) Run() error {
	server := &HTTPServer{
		app:      app,
		services: app.services,
	}
	return server.serveHTTP()
}
