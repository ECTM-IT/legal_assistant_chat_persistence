package db

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/libs"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services/mappers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/mongo"
)

type Services struct {
	AgentService         *services.AgentServiceImpl
	CaseService          *services.CaseServiceImpl
	TeamService          *services.TeamServiceImpl
	UserService          *services.UserServiceImpl
	SubscriptionService  *services.SubscriptionServiceImpl
	PlanService          *services.PlanServiceImpl
	HelpService          *services.HelpServiceImpl
	WebhookService       *services.WebhookServiceImpl
	PaymentMethodService services.PaymentMethodService
	// TODO: uncomment after mailing is up
	// MailerService       libs.MailerService
}

func InitializeServices(db *mongo.Database, logger logs.Logger) *Services {
	// Initialize DAOs
	agentDAO := daos.NewAgentDAO(db, logger)
	caseDAO := daos.NewCaseDAO(db, logger)
	teamDAO := daos.NewTeamDAO(db, logger)
	userDAO := daos.NewUserDAO(db, logger)
	subscriptionDAO := daos.NewSubscriptionsDAO(db, logger)
	invitationDAO := daos.NewInvitationDAO(db, logger)
	stripeEventDAO := daos.NewStripeEventDAO(db, logger)
	paymentMethodDAO := daos.NewPaymentMethodDAO(db, logger)

	// Initialize repositories
	agentRepo := repositories.NewAgentRepository(agentDAO, userDAO, logger)
	caseRepo := repositories.NewCaseRepository(caseDAO)
	teamRepo := repositories.NewTeamRepository(teamDAO, userDAO, invitationDAO, logger)
	userRepo := repositories.NewUserRepository(userDAO)
	subscriptionRepo := repositories.NewSubscriptionRepository(subscriptionDAO)
	stripeEventRepo := repositories.NewStripeEventRepository(stripeEventDAO)
	paymentMethodRepo := repositories.NewPaymentMethodRepository(paymentMethodDAO)

	//Initialize mappers
	agentMapper := mappers.NewAgentConversionService(logger)
	caseMapper := mappers.NewCaseConversionService(logger)
	teamMapper := mappers.NewTeamConversionService(logger)
	userMapper := mappers.NewUserConversionService(logger)
	subscriptionMapper := mappers.NewSubscriptionConversionService(logger)
	planMapper := mappers.NewPlanConversionService(logger)
	paymentMethodMapper := mappers.NewPaymentMethodConversionService(logger)

	// Initialize services
	stripeService := libs.NewStripeService(logger)
	// Cast the stripe service to the extended interface
	updatedStripeService := stripeService.(libs.UpdatedStripeService)

	agentService := services.NewAgentService(agentRepo, agentMapper, userMapper, logger)
	caseService := services.NewCaseService(caseRepo, caseMapper, userMapper, userRepo, logger)
	userService := services.NewUserService(userRepo, subscriptionRepo, userMapper, logger)
	teamService := services.NewTeamService(teamRepo, teamMapper, userRepo, logger)
	planService := services.NewPlanService(subscriptionRepo, planMapper, subscriptionMapper, logger)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo, userRepo, subscriptionMapper, planService, stripeService, logger)
	helpService := services.NewHelpService(nil, logger)
	webhookService := services.NewWebhookService(stripeService, stripeEventRepo, subscriptionRepo, userRepo, logger)
	paymentMethodService := services.NewPaymentMethodService(paymentMethodRepo, userRepo, paymentMethodMapper, updatedStripeService, logger)

	return &Services{
		AgentService:         agentService,
		CaseService:          caseService,
		TeamService:          teamService,
		UserService:          userService,
		SubscriptionService:  subscriptionService,
		PlanService:          planService,
		HelpService:          helpService,
		WebhookService:       webhookService,
		PaymentMethodService: paymentMethodService,
		// TODO: uncomment after mailing is up
		// MailerService:       mailerService,
	}
}
