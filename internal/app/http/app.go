package http

import (
	"os"
	"runtime/debug"
	"sync"

	"go.uber.org/zap"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/utils/env"
)

func Main() {
	// Initialize your ZapLogger
	logger := logs.Init() // Assuming Init() returns an instance of your ZapLogger which implements the Logger interface

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error("Application failed", err, zap.String("trace", trace))
		os.Exit(1)
	}
}

type config struct {
	baseURL  string
	httpPort int
	cookie   struct {
		secretKey string
	}
	// notifications struct { mailing service WIP
	// 	email string
	// }
	// smtp struct {
	// 	host     string
	// 	port     int
	// 	username string
	// 	password string
	// 	from     string
	// }
}

type Application struct {
	config config
	logger logs.Logger // Use your Logger interface here
	wg     sync.WaitGroup
}

func run(logger logs.Logger) error {
	var cfg config

	cfg.baseURL = env.GetString("BASE_URL", "http://localhost:4444")
	cfg.httpPort = env.GetInt("HTTP_PORT", 4444)
	cfg.cookie.secretKey = env.GetString("COOKIE_SECRET_KEY", "3iepwbkq5chsrusjoha26mnsjt233ujq")

	app := &Application{
		config: cfg,
		logger: logger,
	}

	return app.serveHTTP()
}
