// handlers/team_handler.go

package handlers

import (
	"net/http"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
)

type TeamHandler struct {
	service *services.TeamServiceImpl
	BaseHandler
}

func NewTeamHandler(service *services.TeamServiceImpl) *TeamHandler {
	return &TeamHandler{
		service: service,
	}
}

func (h *TeamHandler) CreateTeam(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateTeamRequest
	if err := h.DecodeJSONBody(r, &req); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	team, err := h.service.CreateTeam(r.Context(), req)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to create team")
		return
	}

	h.RespondWithJSON(w, http.StatusCreated, team)
}

func (h *TeamHandler) GetTeamByID(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	team, err := h.service.GetTeamByID(r.Context(), id)
	if err != nil {
		h.RespondWithError(w, http.StatusNotFound, "Team not found")
		return
	}

	h.RespondWithJSON(w, http.StatusOK, team)
}

func (h *TeamHandler) GetAllTeams(w http.ResponseWriter, r *http.Request) {
	teams, err := h.service.GetAllTeams(r.Context())
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve teams")
		return
	}

	h.RespondWithJSON(w, http.StatusOK, teams)
}

func (h *TeamHandler) UpdateTeam(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	var req dtos.UpdateTeamRequest
	if err := h.DecodeJSONBody(r, &req); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	team, err := h.service.UpdateTeam(r.Context(), id, req)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to update team")
		return
	}

	h.RespondWithJSON(w, http.StatusOK, team)
}

func (h *TeamHandler) SoftDeleteTeam(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	err = h.service.SoftDeleteTeam(r.Context(), id)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to delete team")
		return
	}

	h.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Team deleted successfully"})
}

func (h *TeamHandler) UndoTeamDeletion(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	err = h.service.UndoTeamDeletion(r.Context(), id)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to undo team deletion")
		return
	}

	h.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Team deletion undone successfully"})
}

func (h *TeamHandler) AddTeamMember(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	var req dtos.AddTeamMemberRequest
	if err := h.DecodeJSONBody(r, &req); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = h.service.AddTeamMember(r.Context(), id, req)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to add team member")
		return
	}

	h.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Team member added successfully"})
}

func (h *TeamHandler) UpdateTeamMember(w http.ResponseWriter, r *http.Request) {
	teamID, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	memberID, err := h.ParseObjectID(r, "memberId", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid member ID")
		return
	}

	var req dtos.UpdateTeamMemberRequest
	if err := h.DecodeJSONBody(r, &req); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = h.service.UpdateTeamMember(r.Context(), teamID, memberID, req)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to update team member")
		return
	}

	h.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Team member updated successfully"})
}

func (h *TeamHandler) SoftDeleteTeamMember(w http.ResponseWriter, r *http.Request) {
	teamID, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	memberID, err := h.ParseObjectID(r, "memberId", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid member ID")
		return
	}

	err = h.service.SoftDeleteTeamMember(r.Context(), teamID, memberID)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to delete team member")
		return
	}

	h.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Team member deleted successfully"})
}

func (h *TeamHandler) UndoTeamMemberDeletion(w http.ResponseWriter, r *http.Request) {
	teamID, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	memberID, err := h.ParseObjectID(r, "memberId", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid member ID")
		return
	}

	err = h.service.UndoTeamMemberDeletion(r.Context(), teamID, memberID)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to undo team member deletion")
		return
	}

	h.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Team member deletion undone successfully"})
}

func (h *TeamHandler) CreateInvitation(w http.ResponseWriter, r *http.Request) {
	teamID, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	var req dtos.TeamInvitationRequest
	if err := h.DecodeJSONBody(r, &req); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	invitation, err := h.service.CreateInvitation(r.Context(), teamID, req)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to create invitation")
		return
	}

	h.RespondWithJSON(w, http.StatusCreated, invitation)
}

func (h *TeamHandler) AcceptInvitation(w http.ResponseWriter, r *http.Request) {
	var req dtos.AcceptInvitationRequest
	if err := h.DecodeJSONBody(r, &req); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := h.service.AcceptInvitation(r.Context(), req)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to accept invitation")
		return
	}

	h.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Invitation accepted successfully"})
}
