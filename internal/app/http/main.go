package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	handler "github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/handlers"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

func (app *application) serveHTTP() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.httpPort),
		Handler:      handler.Routes(),
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

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		app.logger.Error("Server error", err)
		return err
	}

	err = <-shutdownErrorChan
	if err != nil {
		app.logger.Error("Shutdown error", err)
		return err
	}

	app.logger.Info("Stopped server", zap.String("addr", srv.Addr))

	app.wg.Wait()
	return nil
}
