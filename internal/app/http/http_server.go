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
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/db"
	"go.uber.org/zap"
)

const (
	defaultIdleTimeout    = 24 * time.Hour
	defaultReadTimeout    = 24 * time.Hour
	defaultWriteTimeout   = 24 * time.Hour
	defaultShutdownPeriod = 30 * time.Second
)

type HTTPServer struct {
	app      *Application
	services *db.Services
}

func (s *HTTPServer) serveHTTP() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", s.app.config.HTTPPort),
		Handler:      handler.Routes(s.services.AgentService, s.services.CaseService, s.services.TeamService, s.services.UserService),
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
		s.app.logger.Info("Shutting down server", zap.String("addr", srv.Addr))
		shutdownErrorChan <- srv.Shutdown(ctx)
	}()

	s.app.logger.Info("Starting server", zap.String("addr", srv.Addr))
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "read: connection reset by peer" {
			s.app.logger.Warn("Connection reset by peer" + err.Error())
		} else {
			s.app.logger.Warn("Server error" + err.Error())
			return err
		}
	}

	err = <-shutdownErrorChan
	if err != nil {
		s.app.logger.Warn("Shutdown error" + err.Error())
		return err
	}

	s.app.logger.Info("Stopped server", zap.String("addr", srv.Addr))
	return nil
}
