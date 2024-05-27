package db

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/mongo"
)

type Services struct {
	AgentService        *services.AgentServiceImpl
	CaseService         *services.CaseServiceImpl
	TeamService         *services.TeamServiceImpl
	UserService         *services.UserServiceImpl
	SubscriptionService *services.SubscriptionServiceImpl
}

func InitializeServices(db *mongo.Database, logger logs.Logger) *Services {
	// Initialize DAOs
	agentDAO := daos.NewAgentDAO(db, logger)
	caseDAO := daos.NewCaseDAO(db, logger)
	teamDAO := daos.NewTeamDAO(db, logger)
	userDAO := daos.NewUserDAO(db, logger)
	subscriptionDAO := daos.NewSubscriptionsDAO(db, logger)

	// Initialize repositories
	agentRepo := repositories.NewAgentRepository(agentDAO, userDAO)
	caseRepo := repositories.NewCaseRepository(caseDAO)
	teamRepo := repositories.NewTeamRepository(teamDAO, userDAO)
	userRepo := repositories.NewUserRepository(userDAO)
	subscriptionRepo := repositories.NewSubscriptionRepository(subscriptionDAO)

	// Initialize services
	agentService := services.NewAgentService(agentRepo)
	caseService := services.NewCaseService(caseRepo)
	teamService := services.NewTeamService(teamRepo)
	userService := services.NewUserService(userRepo)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo)

	return &Services{
		AgentService:        agentService,
		CaseService:         caseService,
		TeamService:         teamService,
		UserService:         userService,
		SubscriptionService: subscriptionService,
	}
}
