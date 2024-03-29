package db

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
	"go.mongodb.org/mongo-driver/mongo"
)

type Services struct {
	AgentService *services.AgentService
	CaseService  *services.CaseService
	TeamService  *services.TeamService
	UserService  *services.UserServiceImpl
}

func InitializeServices(db *mongo.Database) *Services {
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

	return &Services{
		AgentService: agentService,
		CaseService:  caseService,
		TeamService:  teamService,
		UserService:  userService,
	}
}
