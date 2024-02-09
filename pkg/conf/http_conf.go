package main

import (
	"errors"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/utils/env"
)

type config struct {
	baseURL  string
	httpPort int
	cookie   struct {
		secretKey string
	}
	notifications struct {
		email string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		from     string
	}
}

func loadConfig() (config, error) {
	cfg := config{
		baseURL:  env.GetString("BASE_URL", "http://localhost:4444"),
		httpPort: env.GetInt("HTTP_PORT", 4444),
		cookie: struct{ secretKey string }{
			secretKey: env.GetString("COOKIE_SECRET_KEY", "3iepwbkq5chsrusjoha26mnsjt233ujq"),
		},
		notifications: struct{ email string }{
			email: env.GetString("NOTIFICATIONS_EMAIL", "default@example.com"),
		},
	}

	if cfg.httpPort < 1024 || cfg.httpPort > 65535 {
		return config{}, errors.New("Invalid value for HTTP_PORT")
	}

	return cfg, nil
}
