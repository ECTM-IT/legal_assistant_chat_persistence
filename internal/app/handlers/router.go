package handler

import (
	"net/http"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/handlers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
	handler "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func Routes(agentService *services.AgentService, caseService *services.CaseService, teamService *services.TeamService, userService *services.UserServiceImpl) http.Handler {
	router := mux.NewRouter()

	headersOk := handler.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handler.AllowedOrigins([]string{"*"})
	methodsOk := handler.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	agentHandler := handlers.NewAgentHandler(agentService)
	caseHandler := handlers.NewCaseHandler(caseService)
	teamHandler := handlers.NewTeamHandler(teamService)
	userHandler := handlers.NewUserHandler(userService)

	// Agent routes
	router.HandleFunc("/agents", agentHandler.GetAllAgents).Methods(http.MethodGet)
	router.HandleFunc("/agents/{id}", agentHandler.GetAgentByID).Methods(http.MethodGet)
	router.HandleFunc("/agents-user", agentHandler.GetAgentsByUser).Methods(http.MethodGet)
	router.HandleFunc("/agent-purchase/{id}", agentHandler.PurchaseAgent).Methods(http.MethodGet)

	// Case routes
	router.HandleFunc("/api-cases", caseHandler.GetAllCases).Methods(http.MethodGet)
	router.HandleFunc("/cases-user", caseHandler.GetCasesByCreatorID).Methods(http.MethodGet)
	router.HandleFunc("/cases/{id}", caseHandler.GetCaseByID).Methods(http.MethodGet)
	router.HandleFunc("/cases-create", caseHandler.CreateCase).Methods(http.MethodPost)
	router.HandleFunc("/cases/{id}", caseHandler.UpdateCase).Methods(http.MethodPut)
	router.HandleFunc("/cases/{id}", caseHandler.DeleteCase).Methods(http.MethodDelete)
	router.HandleFunc("/case-add-user/{id}", caseHandler.AddCollaboratorToCase).Methods(http.MethodPost)
	router.HandleFunc("/case-remove-user/{id}/{userID}", caseHandler.RemoveCollaboratorFromCase).Methods(http.MethodDelete)

	// Team routes
	router.HandleFunc("/teams/{id}", teamHandler.GetTeamByID).Methods(http.MethodGet)
	router.HandleFunc("/team-member/{id}", teamHandler.GetTeamMember).Methods(http.MethodGet)
	router.HandleFunc("/team/change-admin/{id}", teamHandler.ChangeAdmin).Methods(http.MethodPut)
	router.HandleFunc("/team/add/{id}", teamHandler.AddMember).Methods(http.MethodPost)
	router.HandleFunc("/team/remove/{id}/{memberId}", teamHandler.RemoveMember).Methods(http.MethodDelete)

	// User routes
	router.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods(http.MethodGet)
	router.HandleFunc("/users", userHandler.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods(http.MethodPut)
	router.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods(http.MethodDelete)

	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(MethodNotAllowedHandler)

	corsHandler := handler.CORS(originsOk, headersOk, methodsOk)(router)

	return corsHandler
}
