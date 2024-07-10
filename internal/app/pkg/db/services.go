package db

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services/mappers"
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
	agentDAO := daos.NewAgentDAO(db)
	caseDAO := daos.NewCaseDAO(db)
	teamDAO := daos.NewTeamDAO(db)
	userDAO := daos.NewUserDAO(db)
	subscriptionDAO := daos.NewSubscriptionsDAO(db)

	// Initialize repositories
	agentRepo := repositories.NewAgentRepository(agentDAO, userDAO)
	caseRepo := repositories.NewCaseRepository(caseDAO)
	teamRepo := repositories.NewTeamRepository(teamDAO, userDAO)
	userRepo := repositories.NewUserRepository(userDAO)
	subscriptionRepo := repositories.NewSubscriptionRepository(subscriptionDAO)

	//Initialize mappers
	agentMapper := mappers.NewAgentConversionService()
	caseMapper := mappers.NewCaseConversionService()
	teamMapper := mappers.NewTeamConversionService()
	userMapper := mappers.NewUserConversionService()
	subscriptionMapper := mappers.NewSubscriptionConversionService()

	// Initialize services
	agentService := services.NewAgentService(agentRepo, agentMapper, userMapper)
	caseService := services.NewCaseService(caseRepo, caseMapper, userRepo)
	teamService := services.NewTeamService(teamRepo, teamMapper)
	userService := services.NewUserService(userRepo, userMapper)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo, subscriptionMapper)

	return &Services{
		AgentService:        agentService,
		CaseService:         caseService,
		TeamService:         teamService,
		UserService:         userService,
		SubscriptionService: subscriptionService,
	}
}
