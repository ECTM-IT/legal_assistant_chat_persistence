package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
)

type HelpHandler struct {
	helpService *services.HelpServiceImpl
}

func NewHelpHandler(helpService *services.HelpServiceImpl) *HelpHandler {
	return &HelpHandler{
		helpService: helpService,
	}
}

func (h *HelpHandler) SendHelpRequest(w http.ResponseWriter, r *http.Request) {

	var request dtos.HelpRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.helpService.SendHelpEmailFromUser(r.Context(), &request)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Help request sent successfully",
	})
}
