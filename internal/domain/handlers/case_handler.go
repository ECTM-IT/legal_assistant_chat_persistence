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
	caseService *services.CaseService
}

func NewCaseHandler(caseService *services.CaseService) *CaseHandler {
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
	id := strings.TrimSpace(mux.Vars(r)["id"])
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
	println(creatorID)
	cases, err := h.caseService.GetCasesByCreatorID(ctx, creatorID)
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
	err := json.NewDecoder(r.Body).Decode(&caseRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
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
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
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
	updatedCase, err := h.caseService.UpdateCase(ctx, id, updateRequest)
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
	id := strings.TrimSpace(mux.Vars(r)["id"])
	err := h.caseService.DeleteCase(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CaseHandler) AddCollaboratorToCase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := strings.TrimSpace(mux.Vars(r)["id"])
	collaboratorID := mux.Vars(r)["collaboratorID"]
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
	id := strings.TrimSpace(mux.Vars(r)["id"])
	collaboratorID := mux.Vars(r)["collaboratorID"]
	updatedCase, err := h.caseService.RemoveCollaboratorFromCase(ctx, id, collaboratorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCase)
}
