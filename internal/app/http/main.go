package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	handler "github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/handlers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/db"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

const (
	defaultIdleTimeout    = time.Minute
	defaultReadTimeout    = 5 * time.Second
	defaultWriteTimeout   = 10 * time.Second
	defaultShutdownPeriod = 30 * time.Second
)

func bootstrapApplication(db *mongo.Database) (*services.AgentService, *services.CaseService, *services.TeamService, *services.UserServiceImpl) {

	// Initialize DAOs
	agentDAO := daos.NewAgentDAO(db)
	caseDAO := daos.NewCaseDAO(db)
	teamDAO := daos.NewTeamDAO(db)
	userDAO := daos.NewUserDAO(db)

	// Initialize repositories
	agentRepo := repositories.NewAgentRepository(agentDAO, userDAO)
	caseRepo := repositories.NewCaseRepository(caseDAO)
	teamRepo := repositories.NewTeamRepository(teamDAO, userDAO)
	userRepo := repositories.NewUserRepository(userDAO)

	// Initialize services
	agentService := services.NewAgentService(agentRepo)
	caseService := services.NewCaseService(caseRepo)
	teamService := services.NewTeamService(teamRepo)
	userService := services.NewUserService(userRepo)

	return agentService, caseService, teamService, userService
}

func (app *Application) serveHTTP() error {
	uri := "mongodb://localhost:27017/"
	client, err := db.Connect(uri, 60)

	if err != nil {
		return err
	}

	laDatabase := db.CreateDB(client)

	agentService, caseService, teamService, userService := bootstrapApplication(laDatabase)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.httpPort),
		Handler:      handler.Routes(agentService, caseService, teamService, userService),
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
