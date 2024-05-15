package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AgentHandler struct {
	agentService *services.AgentServiceImpl
}

func NewAgentHandler(agentService *services.AgentServiceImpl) *AgentHandler {
	return &AgentHandler{
		agentService: agentService,
	}
}

func (h *AgentHandler) GetAllAgents(w http.ResponseWriter, r *http.Request) {
	agents, err := h.agentService.GetAllAgents(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agents)
}

func (h *AgentHandler) GetAgentByID(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["id"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	agent, err := h.agentService.GetAgentByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agent)
}

func (h *AgentHandler) GetAgentsByUser(w http.ResponseWriter, r *http.Request) {
	userID, err := primitive.ObjectIDFromHex(strings.TrimSpace(r.Header.Get("Authorization")))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	agents, err := h.agentService.GetAgentsByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(agents)
}

func (h *AgentHandler) PurchaseAgent(w http.ResponseWriter, r *http.Request) {
	userID, err := primitive.ObjectIDFromHex(strings.TrimSpace(r.Header.Get("Authorization")))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	agentID, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	user, err := h.agentService.PurchaseAgent(r.Context(), userID, agentID)
	if err != nil {
		if err.Error() == "agent already added to the user" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Agent added successfully",
		"user":    user,
	})
}
