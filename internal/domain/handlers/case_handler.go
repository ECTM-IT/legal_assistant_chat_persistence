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

type CaseHandler struct {
	caseService *services.CaseServiceImpl
}

func NewCaseHandler(caseService *services.CaseServiceImpl) *CaseHandler {
	return &CaseHandler{
		caseService: caseService,
	}
}

func (h *CaseHandler) GetAllCases(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cases, err := h.caseService.GetAllCases(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cases)
}

func (h *CaseHandler) GetCaseByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["id"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	caseResponse, err := h.caseService.GetCaseByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(caseResponse)
}

func (h *CaseHandler) GetCasesByCreatorID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	creatorID := strings.TrimSpace(r.Header.Get("Authorization"))
	objectID, err := primitive.ObjectIDFromHex(creatorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	cases, err := h.caseService.GetCasesByCreatorID(ctx, objectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cases)
}

func (h *CaseHandler) CreateCase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var caseRequest dtos.CreateCaseRequest
	creatorId, err := primitive.ObjectIDFromHex(strings.TrimSpace(r.Header.Get("Authorization")))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&caseRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	caseRequest.CreatorID.Val = creatorId
	createdCase, err := h.caseService.CreateCase(ctx, caseRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdCase)
}

func (h *CaseHandler) UpdateCase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["id"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var updateRequest dtos.UpdateCaseRequest
	err = json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a map to store the update fields
	updateFields := make(map[string]interface{})

	// Check each field in the UpdateCaseRequest and add non-null values to the updateFields map
	if updateRequest.Name.Valid {
		updateFields["name"] = updateRequest.Name.Val
	}
	if updateRequest.Description.Valid {
		updateFields["description"] = updateRequest.Description.Val
	}
	if updateRequest.AgentID.Valid {
		updateFields["agent_id"] = updateRequest.AgentID.Val
	}
	if updateRequest.CollaboratorIDs.Valid {
		updateFields["collaborator_ids"] = updateRequest.CollaboratorIDs.Val
	}
	if updateRequest.Action.Valid {
		updateFields["action"] = updateRequest.Action.Val
	}
	if updateRequest.Skill.Valid {
		updateFields["skill"] = updateRequest.Skill.Val
	}
	if updateRequest.Share.Valid {
		updateFields["share"] = updateRequest.Share.Val
	}
	if updateRequest.IsArchived.Valid {
		updateFields["is_archived"] = updateRequest.IsArchived.Val
	}

	// Call the UpdateCase service method with the update fields
	updatedCase, err := h.caseService.UpdateCase(ctx, id, updateFields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCase)
}

func (h *CaseHandler) DeleteCase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["id"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	deletedCase, err := h.caseService.DeleteCase(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deletedCase)
	w.WriteHeader(http.StatusNoContent)
}

func (h *CaseHandler) AddCollaboratorToCase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["id"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	collaboratorID, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["collaboratorID"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	updatedCase, err := h.caseService.AddCollaboratorToCase(ctx, id, collaboratorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCase)
}

func (h *CaseHandler) RemoveCollaboratorFromCase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["id"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	collaboratorID, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["collaboratorID"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	updatedCase, err := h.caseService.RemoveCollaboratorFromCase(ctx, id, collaboratorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCase)
}
