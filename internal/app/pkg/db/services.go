package db

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/security" // Importing the security package
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
	EncryptionService   *security.EncryptionService
}

func InitializeServices(db *mongo.Database, logger logs.Logger) *Services {

	// Initialize the EnvironmentKeyManager
	keyManager := security.NewEnvironmentKeyManager("ENV_AES_256_KEY")

	// Initialize EncryptionService
	encryptionService, err := security.NewAES256EncryptionService(keyManager)
	if err != nil {
		logger.Error("error in encryption service", err)
	}
	// Initialize DAOs
	agentDAO := daos.NewAgentDAO(db, logger)
	caseDAO := daos.NewCaseDAO(db, logger)
	teamDAO := daos.NewTeamDAO(db, logger)
	userDAO := daos.NewUserDAO(db, logger)
	subscriptionDAO := daos.NewSubscriptionsDAO(db, logger)

	// Initialize repositories
	agentRepo := repositories.NewAgentRepository(agentDAO, userDAO, logger)
	caseRepo := repositories.NewCaseRepository(caseDAO)
	teamRepo := repositories.NewTeamRepository(teamDAO, userDAO, logger)
	userRepo := repositories.NewUserRepository(userDAO)
	subscriptionRepo := repositories.NewSubscriptionRepository(subscriptionDAO)

	//Initialize mappers
	agentMapper := mappers.NewAgentConversionService(logger)
	caseMapper := mappers.NewCaseConversionService(logger)
	teamMapper := mappers.NewTeamConversionService(logger)
	userMapper := mappers.NewUserConversionService(logger)
	subscriptionMapper := mappers.NewSubscriptionConversionService(logger)

	// Initialize services
	agentService := services.NewAgentService(agentRepo, agentMapper, userMapper, logger)
	caseService := services.NewCaseService(caseRepo, caseMapper, userMapper, userRepo, logger)
	teamService := services.NewTeamService(teamRepo, teamMapper, logger)
	userService := services.NewUserService(userRepo, userMapper, logger, encryptionService) // Pass encryptionService to UserService
	subscriptionService := services.NewSubscriptionService(subscriptionRepo, subscriptionMapper, logger)

	return &Services{
		AgentService:        agentService,
		CaseService:         caseService,
		TeamService:         teamService,
		UserService:         userService,
		SubscriptionService: subscriptionService,
	}
}
