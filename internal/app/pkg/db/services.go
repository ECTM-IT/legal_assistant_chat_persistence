package db

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/config"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/libs"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services/mappers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type Services struct {
	AgentService        *services.AgentServiceImpl
	CaseService         *services.CaseServiceImpl
	TeamService         *services.TeamServiceImpl
	UserService         *services.UserServiceImpl
	SubscriptionService *services.SubscriptionServiceImpl
	PlanService         *services.PlanServiceImpl
	HelpService         *services.HelpServiceImpl
	MailerService       libs.MailerService
}

func InitializeServices(db *mongo.Database, gcpSaKeyPath string, logger logs.Logger) *Services {
	// Create a context
	ctx := context.Background()

	// Initialize DAOs
	agentDAO := daos.NewAgentDAO(db, logger)
	caseDAO := daos.NewCaseDAO(db, logger)
	teamDAO := daos.NewTeamDAO(db, logger)
	userDAO := daos.NewUserDAO(db, logger)
	subscriptionDAO := daos.NewSubscriptionsDAO(db, logger)
	invitationDAO := daos.NewInvitationDAO(db, logger)

	// Initialize repositories
	agentRepo := repositories.NewAgentRepository(agentDAO, userDAO, logger)
	caseRepo := repositories.NewCaseRepository(caseDAO)
	teamRepo := repositories.NewTeamRepository(teamDAO, userDAO, invitationDAO, logger)
	userRepo := repositories.NewUserRepository(userDAO)
	subscriptionRepo := repositories.NewSubscriptionRepository(subscriptionDAO)

	//Initialize mappers
	agentMapper := mappers.NewAgentConversionService(logger)
	caseMapper := mappers.NewCaseConversionService(logger)
	teamMapper := mappers.NewTeamConversionService(logger)
	userMapper := mappers.NewUserConversionService(logger)
	subscriptionMapper := mappers.NewSubscriptionConversionService(logger)
	planMapper := mappers.NewPlanConversionService(logger)

	// Load mailer configuration
	mailerConfig := config.LoadMailerConfig()

	// Authenticate using the service account
	srv, err := drive.NewService(ctx, option.WithCredentialsFile(gcpSaKeyPath))
	if err != nil {
		logger.Error("Unable to create Drive client: %v", err)
	}

	// initialize Google Drive Service
	driveService := services.NewDriveService(srv, logger)

	// Initialize services
	mailerService := libs.NewMailerService(
		mailerConfig.Host,
		mailerConfig.Port,
		mailerConfig.Username,
		mailerConfig.Password,
		mailerConfig.From,
		logger,
	)
	agentService := services.NewAgentService(agentRepo, agentMapper, userMapper, logger)
	caseService := services.NewCaseService(caseRepo, caseMapper, userMapper, userRepo, driveService, logger)
	teamService := services.NewTeamService(teamRepo, teamMapper, logger)
	userService := services.NewUserService(userRepo, userMapper, logger)
	planService := services.NewPlanService(subscriptionRepo, planMapper, subscriptionMapper, logger)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo, userRepo, subscriptionMapper, planService, mailerService, logger)
	helpService := services.NewHelpService(mailerService, logger)

	return &Services{
		AgentService:        agentService,
		CaseService:         caseService,
		TeamService:         teamService,
		UserService:         userService,
		SubscriptionService: subscriptionService,
		PlanService:         planService,
		HelpService:         helpService,
		MailerService:       mailerService,
	}
}
