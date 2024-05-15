package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeamHandler struct {
	teamService *services.TeamServiceImpl
}

func NewTeamHandler(teamService *services.TeamServiceImpl) *TeamHandler {
	return &TeamHandler{
		teamService: teamService,
	}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var request dtos.CreateTeamRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	team, err := h.teamService.CreateTeam(r.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(team)
}

func (h *TeamHandler) GetTeamByID(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["id"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	team, err := h.teamService.GetTeamByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

func (h *TeamHandler) GetAllTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := h.teamService.GetAllTeams(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teams)
}

func (h *TeamHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["id"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	var request dtos.UpdateTeamRequest
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	team, err := h.teamService.UpdateTeam(r.Context(), id, request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(team)
}

func (h *TeamHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["id"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	err = h.teamService.DeleteTeam(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TeamHandler) GetTeamMember(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["id"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	member, err := h.teamService.GetTeamMember(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

func (h *TeamHandler) ChangeAdmin(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["id"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	var request dtos.ChangeAdminRequest
	err = json.NewDecoder(r.Body).Decode(&request)
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
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["id"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	var request dtos.AddMemberRequest
	err = json.NewDecoder(r.Body).Decode(&request)
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
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["id"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	memberID, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["memberId"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	removedMember, err := h.teamService.RemoveMember(r.Context(), id, memberID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(removedMember)
}
