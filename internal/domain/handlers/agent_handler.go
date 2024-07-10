package handlers

import (
	"context"
	"net/http"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AgentService interface {
	GetAllAgents(ctx context.Context) ([]dtos.AgentResponse, error)
	GetAgentByID(ctx context.Context, id primitive.ObjectID) (*dtos.AgentResponse, error)
	GetAgentsByUserID(ctx context.Context, userID primitive.ObjectID) ([]dtos.AgentResponse, error)
	PurchaseAgent(ctx context.Context, userID, agentID primitive.ObjectID) (*dtos.UserResponse, error)
}

type AgentHandler struct {
	BaseHandler
	service *services.AgentServiceImpl
}

func NewAgentHandler(service *services.AgentServiceImpl) *AgentHandler {
	return &AgentHandler{service: service}
}

func (h *AgentHandler) GetAllAgents(w http.ResponseWriter, r *http.Request) {
	agents, err := h.service.GetAllAgents(r.Context())
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve agents")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, agents)
}

func (h *AgentHandler) GetAgentByID(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid agent ID")
		return
	}

	agent, err := h.service.GetAgentByID(r.Context(), id)
	if err != nil {
		h.RespondWithError(w, http.StatusNotFound, "Agent not found")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, agent)
}

func (h *AgentHandler) GetAgentsByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := h.ParseObjectID(r, "userID", true)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	agents, err := h.service.GetAgentsByUserID(r.Context(), userID)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve agents")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, agents)
}

func (h *AgentHandler) PurchaseAgent(w http.ResponseWriter, r *http.Request) {
	userID, err := h.ParseObjectID(r, "userID", true)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	agentID, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid agent ID")
		return
	}

	user, err := h.service.PurchaseAgent(r.Context(), userID, agentID)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to purchase agent")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, user)
}
