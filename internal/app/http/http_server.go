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
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.uber.org/zap"
)

const (
	defaultIdleTimeout    = 24 * time.Hour
	defaultReadTimeout    = 24 * time.Hour
	defaultWriteTimeout   = 24 * time.Hour
	defaultShutdownPeriod = 30 * time.Second
)

// HTTPServer holds the necessary components to run the HTTP server.
type HTTPServer struct {
	addr       string
	handler    http.Handler
	logger     logs.Logger
	shutdownCh chan os.Signal
}

// NewHTTPServer initializes a new HTTPServer with the provided configuration and dependencies.
func NewHTTPServer(config *Config, services *db.Services, logger logs.Logger) *HTTPServer {
	return &HTTPServer{
		addr:       fmt.Sprintf("0.0.0.0:%d", config.HTTPPort),
		handler:    handler.Routes(services.AgentService, services.CaseService, services.TeamService, services.UserService, services.SubscriptionService, services.PlanService),
		logger:     logger,
		shutdownCh: make(chan os.Signal, 1),
	}
}

// ServeHTTP starts the HTTP server and manages graceful shutdown.
func (s *HTTPServer) ServeHTTP() error {
	srv := &http.Server{
		Addr:         s.addr,
		Handler:      s.handler,
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	// Handle graceful shutdown
	signal.Notify(s.shutdownCh, syscall.SIGINT, syscall.SIGTERM)
	shutdownErrorChan := make(chan error)

	go func() {
		<-s.shutdownCh
		ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownPeriod)
		defer cancel()
		s.logger.Info("Shutting down server", zap.String("addr", srv.Addr))
		shutdownErrorChan <- srv.Shutdown(ctx)
	}()

	s.logger.Info("Starting server", zap.String("addr", srv.Addr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "read: connection reset by peer" {
			s.logger.Warn("Connection reset by peer", zap.Error(err))
		} else {
			s.logger.Error("Server error", err)
			return err
		}
	}

	if err := <-shutdownErrorChan; err != nil {
		s.logger.Error("Shutdown error", err)
		return err
	}

	s.logger.Info("Stopped server", zap.String("addr", srv.Addr))
	return nil
}
