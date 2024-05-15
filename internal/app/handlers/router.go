package handler

import (
	"net/http"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/handlers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Routes initializes the routes for the application with the provided services.
func Routes(agentService *services.AgentServiceImpl, caseService *services.CaseServiceImpl, teamService *services.TeamServiceImpl, userService *services.UserServiceImpl) http.Handler {
	router := mux.NewRouter()

	// Create a new CORS handler with the desired configuration
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
			http.MethodDelete,
			http.MethodPatch,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})

	agentHandler := handlers.NewAgentHandler(agentService)
	caseHandler := handlers.NewCaseHandler(caseService)
	teamHandler := handlers.NewTeamHandler(teamService)
	userHandler := handlers.NewUserHandler(userService)

	// Register agent routes
	registerAgentRoutes(router, agentHandler)

	// Register case routes
	registerCaseRoutes(router, caseHandler)

	// Register team routes
	registerTeamRoutes(router, teamHandler)

	// Register user routes
	registerUserRoutes(router, userHandler)

	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(MethodNotAllowedHandler)

	// Wrap the router with the CORS handler
	return corsHandler.Handler(router)
}

func registerAgentRoutes(router *mux.Router, handler *handlers.AgentHandler) {
	router.HandleFunc("/agents/", handler.GetAllAgents).Methods(http.MethodGet)
	router.HandleFunc("/agents/{id}/", handler.GetAgentByID).Methods(http.MethodGet)
	router.HandleFunc("/agents-user/", handler.GetAgentsByUser).Methods(http.MethodGet)
	router.HandleFunc("/agent-purchase/{id}/", handler.PurchaseAgent).Methods(http.MethodGet)
}

func registerCaseRoutes(router *mux.Router, handler *handlers.CaseHandler) {
	router.HandleFunc("/api-cases/", handler.GetAllCases).Methods(http.MethodGet)
	router.HandleFunc("/cases-user/", handler.GetCasesByCreatorID).Methods(http.MethodGet)
	router.HandleFunc("/cases/{id}/", handler.GetCaseByID).Methods(http.MethodGet)
	router.HandleFunc("/cases-create/", handler.CreateCase).Methods(http.MethodPost)
	router.HandleFunc("/cases/{id}/", handler.UpdateCase).Methods(http.MethodPatch)
	router.HandleFunc("/cases/{id}/", handler.DeleteCase).Methods(http.MethodDelete)
	router.HandleFunc("/case-add-user/{id}/", handler.AddCollaboratorToCase).Methods(http.MethodPost)
	router.HandleFunc("/case-remove-user/{id}/{userID}/", handler.RemoveCollaboratorFromCase).Methods(http.MethodDelete)
}

func registerTeamRoutes(router *mux.Router, handler *handlers.TeamHandler) {
	router.HandleFunc("/teams/{id}/", handler.GetTeamByID).Methods(http.MethodGet)
	router.HandleFunc("/team-member/{id}/", handler.GetTeamMember).Methods(http.MethodGet)
	router.HandleFunc("/team/change-admin/{id}/", handler.ChangeAdmin).Methods(http.MethodPatch)
	router.HandleFunc("/team/add/{id}/", handler.AddMember).Methods(http.MethodPost)
	router.HandleFunc("/team/remove/{id}/{memberId}/", handler.RemoveMember).Methods(http.MethodDelete)
}

func registerUserRoutes(router *mux.Router, handler *handlers.UserHandler) {
	router.HandleFunc("/users/{id}/", handler.GetUserByID).Methods(http.MethodGet)
	router.HandleFunc("/users-email/", handler.GetUserByEmail).Methods(http.MethodPost)
	router.HandleFunc("/users/", handler.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/users/{id}/", handler.UpdateUser).Methods(http.MethodPatch)
	router.HandleFunc("/users/{id}/", handler.DeleteUser).Methods(http.MethodDelete)
}
