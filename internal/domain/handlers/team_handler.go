// handlers/team_handler.go

package handlers

import (
	"context"
	"net/http"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeamService interface {
	CreateTeam(ctx context.Context, req dtos.CreateTeamRequest) (*dtos.TeamResponse, error)
	GetTeamByID(ctx context.Context, id primitive.ObjectID) (*dtos.TeamResponse, error)
	GetAllTeams(ctx context.Context) ([]dtos.TeamResponse, error)
	UpdateTeam(ctx context.Context, id primitive.ObjectID, req dtos.UpdateTeamRequest) (*dtos.TeamResponse, error)
	DeleteTeam(ctx context.Context, id primitive.ObjectID) error
	GetTeamMember(ctx context.Context, id primitive.ObjectID) (*dtos.TeamMemberResponse, error)
	ChangeAdmin(ctx context.Context, id primitive.ObjectID, req dtos.ChangeAdminRequest) (*dtos.TeamMemberResponse, error)
	AddMember(ctx context.Context, id primitive.ObjectID, req dtos.AddMemberRequest) (*dtos.TeamMemberResponse, error)
	RemoveMember(ctx context.Context, id, memberID primitive.ObjectID) (*dtos.TeamMemberResponse, error)
}

type TeamHandler struct {
	BaseHandler
	service *services.TeamServiceImpl
}

func NewTeamHandler(service *services.TeamServiceImpl) *TeamHandler {
	return &TeamHandler{service: service}
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

func (h *TeamHandler) DeleteTeam(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	err = h.service.DeleteTeam(r.Context(), id)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to delete team")
		return
	}
	h.RespondWithJSON(w, http.StatusNoContent, nil)
}

func (h *TeamHandler) GetTeamMember(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	member, err := h.service.GetTeamMember(r.Context(), id)
	if err != nil {
		h.RespondWithError(w, http.StatusNotFound, "Team member not found")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, member)
}

func (h *TeamHandler) ChangeAdmin(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	var req dtos.ChangeAdminRequest
	if err := h.DecodeJSONBody(r, &req); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	member, err := h.service.ChangeAdmin(r.Context(), id, req)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to change admin")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, member)
}

func (h *TeamHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	var req dtos.AddMemberRequest
	if err := h.DecodeJSONBody(r, &req); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	member, err := h.service.AddMember(r.Context(), id, req)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to add member")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, member)
}

func (h *TeamHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid team ID")
		return
	}

	memberID, err := h.ParseObjectID(r, "memberId", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid member ID")
		return
	}

	removedMember, err := h.service.RemoveMember(r.Context(), id, memberID)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to remove member")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, removedMember)
}
