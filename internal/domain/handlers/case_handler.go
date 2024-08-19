package handlers

import (
	"context"
	"net/http"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CaseHanler interface {
	GetAllCases(ctx context.Context) ([]dtos.CaseResponse, error)
	GetCaseByID(ctx context.Context, id primitive.ObjectID) (*dtos.CaseResponse, error)
	GetCasesByCreatorID(ctx context.Context, creatorID primitive.ObjectID) ([]dtos.CaseResponse, error)
	CreateCase(ctx context.Context, req *dtos.CreateCaseRequest) (*dtos.CaseResponse, error)
	UpdateCase(ctx context.Context, id primitive.ObjectID, req *dtos.UpdateCaseRequest) (*dtos.CaseResponse, error)
	DeleteCase(ctx context.Context, id primitive.ObjectID) (*dtos.CaseResponse, error)
	AddCollaboratorToCase(ctx context.Context, caseID, collaboratorID primitive.ObjectID) (*dtos.CaseResponse, error)
	RemoveCollaboratorFromCase(ctx context.Context, caseID, collaboratorID primitive.ObjectID) (*dtos.CaseResponse, error)
}

type CaseHandler struct {
	BaseHandler
	service *services.CaseServiceImpl
}

func NewCaseHandler(service *services.CaseServiceImpl) *CaseHandler {
	return &CaseHandler{service: service}
}

func (h *CaseHandler) GetAllCases(w http.ResponseWriter, r *http.Request) {
	cases, err := h.service.GetAllCases(r.Context())
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve cases")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, cases)
}

func (h *CaseHandler) GetCaseByID(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid case ID")
		return
	}

	caseResponse, err := h.service.GetCaseByID(r.Context(), id)
	if err != nil {
		h.RespondWithError(w, http.StatusNotFound, "Case not found")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, caseResponse)
}

func (h *CaseHandler) GetCasesByCreatorID(w http.ResponseWriter, r *http.Request) {
	creatorID, err := h.ParseObjectID(r, "", true)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid creator ID")
		return
	}

	cases, err := h.service.GetCasesByCreatorID(r.Context(), creatorID)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve cases")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, cases)
}

func (h *CaseHandler) CreateCase(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateCaseRequest
	if err := h.DecodeJSONBody(r, &req); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	creatorID, err := h.BaseHandler.ParseObjectID(r, "Authorization", true)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "CreatorID not found on headers")
	}
	req.CreatorID.Value = creatorID
	req.CreatorID.Present = true
	createdCase, err := h.service.CreateCase(r.Context(), req)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to create case")
		return
	}
	h.RespondWithJSON(w, http.StatusCreated, createdCase)
}

func (h *CaseHandler) UpdateCase(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid case ID")
		return
	}

	var req dtos.UpdateCaseRequest
	if err := h.DecodeJSONBody(r, &req); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedCase, err := h.service.UpdateCase(r.Context(), id, req)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to update case")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, updatedCase)
}

func (h *CaseHandler) DeleteCase(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid case ID")
		return
	}

	deletedCase, err := h.service.DeleteCase(r.Context(), id)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to delete case")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, deletedCase)
}

func (h *CaseHandler) AddCollaboratorToCase(w http.ResponseWriter, r *http.Request) {
	caseID, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid case ID")
		return
	}

	collaboratorID, err := h.ParseObjectID(r, "collaboratorID", true)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid collaborator ID")
		return
	}

	updatedCase, err := h.service.AddCollaboratorToCase(r.Context(), caseID, collaboratorID)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to add collaborator to case")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, updatedCase)
}

func (h *CaseHandler) RemoveCollaboratorFromCase(w http.ResponseWriter, r *http.Request) {
	caseID, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid case ID")
		return
	}

	collaboratorID, err := h.ParseObjectID(r, "collaboratorID", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid collaborator ID")
		return
	}

	updatedCase, err := h.service.RemoveCollaboratorFromCase(r.Context(), caseID, collaboratorID)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to remove collaborator from case")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, updatedCase)
}
