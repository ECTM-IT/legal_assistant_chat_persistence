package handler

import (
	"net/http"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/handlers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
	"github.com/gorilla/mux"
)

// CORS middleware to allow all methods
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type, Authorization")
		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		// Pass to the next handler
		next.ServeHTTP(w, r)
	})
}

// Routes initializes the routes for the application with the provided services.
func Routes(
	agentService *services.AgentServiceImpl,
	caseService *services.CaseServiceImpl,
	teamService *services.TeamServiceImpl,
	userService *services.UserServiceImpl,
	subscriptionService *services.SubscriptionServiceImpl,
	planService *services.PlanServiceImpl,
) http.Handler {
	router := mux.NewRouter()

	// Create a new CORS handler with the desired configuration
	// corsHandler := cors.New(cors.Options{
	// 	AllowedOrigins: []string{"*"},
	// 	AllowedMethods: []string{
	// 		http.MethodPost,
	// 		http.MethodGet,
	// 		http.MethodDelete,
	// 		http.MethodPatch,
	// 		http.MethodOptions,
	// 	},
	// 	AllowedHeaders:   []string{"*"},
	// 	AllowCredentials: false,
	// })

	agentHandler := handlers.NewAgentHandler(agentService)
	caseHandler := handlers.NewCaseHandler(caseService)
	teamHandler := handlers.NewTeamHandler(teamService)
	userHandler := handlers.NewUserHandler(userService)
	subscriptionHandler := handlers.NewSubscriptionHandler(subscriptionService)
	planHandler := handlers.NewPlanHandler(planService)
	// Register agent routes
	registerAgentRoutes(router, agentHandler)

	// Register case routes
	registerCaseRoutes(router, caseHandler)

	// Register team routes
	registerTeamRoutes(router, teamHandler)

	// Register user routes
	registerUserRoutes(router, userHandler)

	// Register subscription routes
	registerSubscriptionRoutes(router, subscriptionHandler)

	// Register plan routes
	registerPlanRoutes(router, planHandler)

	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(MethodNotAllowedHandler)

	// Wrap the router with the CORS handler
	return corsMiddleware(router)
}

func registerAgentRoutes(router *mux.Router, handler *handlers.AgentHandler) {
	router.HandleFunc("/agents/", handler.GetAllAgents).Methods(http.MethodGet)
	router.HandleFunc("/agents/{id}/", handler.GetAgentByID).Methods(http.MethodGet)
	router.HandleFunc("/agents-user/", handler.GetAgentsByUser).Methods(http.MethodGet)
	router.HandleFunc("/agent-purchase/{id}/", handler.PurchaseAgent).Methods(http.MethodGet)
}

func registerTeamRoutes(router *mux.Router, handler *handlers.TeamHandler) {
	router.HandleFunc("/teams/", handler.CreateTeam).Methods(http.MethodPost)
	router.HandleFunc("/teams/{id}/", handler.GetTeamByID).Methods(http.MethodGet)
	router.HandleFunc("/teams/", handler.GetAllTeams).Methods(http.MethodGet)
	router.HandleFunc("/teams/{id}/", handler.UpdateTeam).Methods(http.MethodPut)
	router.HandleFunc("/teams/{id}/", handler.SoftDeleteTeam).Methods(http.MethodDelete)
	router.HandleFunc("/teams/{id}/restore/", handler.UndoTeamDeletion).Methods(http.MethodPost)
	router.HandleFunc("/teams/{id}/members/", handler.AddTeamMember).Methods(http.MethodPost)
	router.HandleFunc("/teams/{id}/members/{memberId}/", handler.UpdateTeamMember).Methods(http.MethodPut)
	router.HandleFunc("/teams/{id}/members/{memberId}/", handler.SoftDeleteTeamMember).Methods(http.MethodDelete)
	router.HandleFunc("/teams/{id}/members/{memberId}/restore/", handler.UndoTeamMemberDeletion).Methods(http.MethodPost)
	router.HandleFunc("/teams/{id}/invitations/", handler.CreateInvitation).Methods(http.MethodPost)
	router.HandleFunc("/teams/invitations/accept/", handler.AcceptInvitation).Methods(http.MethodPost)
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
	router.HandleFunc("/case-add-document/{id}/", handler.AddDocumentToCase).Methods(http.MethodPost)
	router.HandleFunc("/case-update-document/{id}/document/{documentID}/", handler.UpdateDocument).Methods(http.MethodPatch)
	router.HandleFunc("/case-add-document-collaborator/{id}/document/{documentID}/", handler.AddDocumentCollaborator).Methods(http.MethodPost)
	router.HandleFunc("/case-remove-document/{id}/{documentID}/", handler.DeleteDocumentFromCase).Methods(http.MethodDelete)
	router.HandleFunc("/case-add-feedback-to-message/{id}/{messageId}/", handler.AddFeedbackToMessage).Methods(http.MethodPost)
	router.HandleFunc("/case-get-feedback-from-message/{id}/{creatorId}/{messageId}/", handler.GetCaseByID).Methods(http.MethodGet)
}

func registerUserRoutes(router *mux.Router, handler *handlers.UserHandler) {
	router.HandleFunc("/users/{id}/", handler.GetUserByID).Methods(http.MethodGet)
	router.HandleFunc("/users-email/", handler.GetUserByEmail).Methods(http.MethodPost)
	router.HandleFunc("/users/", handler.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/users/{id}/", handler.UpdateUser).Methods(http.MethodPatch)
	router.HandleFunc("/users/{id}/", handler.DeleteUser).Methods(http.MethodDelete)
}

func registerSubscriptionRoutes(router *mux.Router, handler *handlers.SubscriptionHandler) {
	router.HandleFunc("/subscriptions/", handler.GetAllSubscriptions).Methods(http.MethodGet)
	router.HandleFunc("/subscriptions/{id}/", handler.GetSubscriptionByID).Methods(http.MethodGet)
	router.HandleFunc("/subscriptions/plan/", handler.GetSubscriptionsByPlan).Methods(http.MethodGet)
	router.HandleFunc("/subscriptions/", handler.CreateSubscription).Methods(http.MethodPost)
	router.HandleFunc("/subscriptions/{id}/", handler.UpdateSubscription).Methods(http.MethodPatch)
	router.HandleFunc("/subscriptions/{id}/", handler.DeleteSubscription).Methods(http.MethodDelete)
	router.HandleFunc("/subscriptions/purchase/", handler.PurchaseSubscription).Methods(http.MethodPost)
}

func registerPlanRoutes(router *mux.Router, handler *handlers.PlanHandler) {
	router.HandleFunc("/plans/", handler.GetPlanOptions).Methods(http.MethodGet)
	router.HandleFunc("/plans/toggle/", handler.TogglePlanType).Methods(http.MethodPatch)
	router.HandleFunc("/plans/select/", handler.SelectPlan).Methods(http.MethodPost)
}
