package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	handler "github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/handlers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/db"
	"go.uber.org/zap"
)

const (
	defaultIdleTimeout    = 24 * time.Hour
	defaultReadTimeout    = 24 * time.Hour
	defaultWriteTimeout   = 24 * time.Hour
	defaultShutdownPeriod = 30 * time.Second
)

func (app *Application) serveHTTP() error {
	uri := "mongodb://0.0.0.0:27017/"
	client, err := db.Connect(uri, 60)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())

	laDatabase := db.CreateDB(client)

	services := db.InitializeServices(laDatabase)

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", app.config.HTTPPort),
		Handler:      handler.Routes(services.AgentService, services.CaseService, services.TeamService, services.UserService),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	shutdownErrorChan := make(chan error)
	go func() {
		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)
		<-quitChan

		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		defer cancel()

		app.logger.Info("Shutting down server", zap.String("addr", srv.Addr))
		shutdownErrorChan <- srv.Shutdown(ctx)
	}()

	app.logger.Info("Starting server", zap.String("addr", srv.Addr))
	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "read: connection reset by peer" {
			app.logger.Warn("Connection reset by peer" + err.Error())
		} else {
			app.logger.Warn("Server error" + err.Error())
			return err
		}
	}

	err = <-shutdownErrorChan
	if err != nil {
		app.logger.Warn("Shutdown error" + err.Error())
		return err
	}

	app.logger.Info("Stopped server", zap.String("addr", srv.Addr))
	return nil
}
