package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
	"github.com/gorilla/mux"
)

type TeamHandler struct {
	teamService *services.TeamService
}

func NewTeamHandler(teamService *services.TeamService) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

func (h *TeamHandler) GetTeamByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	team, err := h.teamService.GetTeamByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

func (h *TeamHandler) GetTeamMember(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	member, err := h.teamService.GetTeamMember(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

func (h *TeamHandler) ChangeAdmin(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var request dtos.ChangeAdminRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	member, err := h.teamService.ChangeAdmin(r.Context(), id, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

func (h *TeamHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var request dtos.AddMemberRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	member, err := h.teamService.AddMember(r.Context(), id, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

func (h *TeamHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	memberID := mux.Vars(r)["memberId"]

	err := h.teamService.RemoveMember(r.Context(), id, memberID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
