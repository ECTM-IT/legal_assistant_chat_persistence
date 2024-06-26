package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
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
		h.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve cases")
		return
	}
	h.respondWithJSON(w, http.StatusOK, cases)
}

func (h *CaseHandler) GetCaseByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := h.getObjectIDFromVars(r, "id")
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "InValueid case ID")
		return
	}
	caseResponse, err := h.caseService.GetCaseByID(ctx, id)
	if err != nil {
		h.respondWithError(w, http.StatusNotFound, "Case not found")
		return
	}
	h.respondWithJSON(w, http.StatusOK, caseResponse)
}

func (h *CaseHandler) GetCasesByCreatorID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	creatorID, err := h.getObjectIDFromHeader(r, "Authorization")
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "InValueid creator ID")
		return
	}
	cases, err := h.caseService.GetCasesByCreatorID(ctx, creatorID)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to retrieve cases")
		return
	}
	h.respondWithJSON(w, http.StatusOK, cases)
}

func (h *CaseHandler) CreateCase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var caseRequest dtos.CreateCaseRequest
	if err := json.NewDecoder(r.Body).Decode(&caseRequest); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "InValueid request payload")
		return
	}
	creatorID, err := h.getObjectIDFromHeader(r, "Authorization")
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "InValueid creator ID")
		return
	}
	caseRequest.CreatorID.Value = creatorID

	caseModel, err := h.prepareCaseModel(caseRequest)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	createdCase, err := h.caseService.CreateCase(ctx, *caseModel)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create case")
		return
	}
	h.respondWithJSON(w, http.StatusCreated, createdCase)
}

func (h *CaseHandler) UpdateCase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := h.getObjectIDFromVars(r, "id")
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "InValueid case ID")
		return
	}

	var updateRequest dtos.UpdateCaseRequest
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "InValueid request payload")
		return
	}

	updateFields := h.prepareUpdateFields(updateRequest)
	updatedCase, err := h.caseService.UpdateCase(ctx, id, updateFields)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to update case")
		return
	}
	h.respondWithJSON(w, http.StatusOK, updatedCase)
}

func (h *CaseHandler) DeleteCase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id, err := h.getObjectIDFromVars(r, "id")
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "InValueid case ID")
		return
	}
	deletedCase, err := h.caseService.DeleteCase(ctx, id)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to delete case")
		return
	}
	h.respondWithJSON(w, http.StatusOK, deletedCase)
}

func (h *CaseHandler) AddCollaboratorToCase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	caseID, collaboratorID, err := h.getCaseAndCollaboratorIDs(r)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "InValueid case or collaborator ID")
		return
	}
	updatedCase, err := h.caseService.AddCollaboratorToCase(ctx, caseID, collaboratorID)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to add collaborator")
		return
	}
	h.respondWithJSON(w, http.StatusOK, updatedCase)
}

func (h *CaseHandler) RemoveCollaboratorFromCase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	caseID, collaboratorID, err := h.getCaseAndCollaboratorIDs(r)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "InValueid case or collaborator ID")
		return
	}
	updatedCase, err := h.caseService.RemoveCollaboratorFromCase(ctx, caseID, collaboratorID)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to remove collaborator")
		return
	}
	h.respondWithJSON(w, http.StatusOK, updatedCase)
}

// Helper methods

func (h *CaseHandler) prepareCaseModel(caseRequest dtos.CreateCaseRequest) (*models.Case, error) {
	messages := h.prepareMessages(caseRequest.Messages.Value)
	collaborators := h.prepareCollaborators(caseRequest)

	return &models.Case{
		ID:            primitive.NewObjectID(),
		Name:          caseRequest.Name.OrElse("New Case"),
		Description:   caseRequest.Description.OrElse(""),
		CreatorID:     caseRequest.CreatorID.Value,
		Messages:      messages,
		Collaborators: collaborators,
		Action:        caseRequest.Action.OrElse("Riassumere"),
		AgentID:       caseRequest.AgentID.OrElse(primitive.NilObjectID),
		LastEdit:      caseRequest.LastEdit.OrElse(time.Now()),
		Share:         caseRequest.Share.OrElse(false),
		IsArchived:    caseRequest.IsArchived.OrElse(false),
	}, nil
}

func (h *CaseHandler) prepareMessages(messagesDTO []dtos.MessageResponse) []models.Message {

	messages := make([]models.Message, len(messagesDTO))
	for i, msg := range messagesDTO {
		messages[i] = models.Message{
			Content:     msg.Content.OrElse(""),
			SenderID:    msg.Sender.OrElse(""),
			RecipientID: msg.Recipient.OrElse(""),
			Skill:       msg.Skill.OrElse(""),
		}
	}
	return messages
}

func (h *CaseHandler) prepareCollaborators(caseRequest dtos.CreateCaseRequest) []models.Collaborators {
	collaborators := []models.Collaborators{{
		ID:   caseRequest.CreatorID.Value,
		Edit: true,
	}}

	if caseRequest.Collaborators.Present {
		for _, collab := range caseRequest.Collaborators.Value {
			collaborators = append(collaborators, models.Collaborators{
				ID:   collab.ID.OrElse(primitive.NilObjectID),
				Edit: collab.Edit.OrElse(false),
			})
		}
	}

	return collaborators
}

func (h *CaseHandler) prepareUpdateFields(updateRequest dtos.UpdateCaseRequest) map[string]interface{} {
	updateFields := make(map[string]interface{})

	if updateRequest.Messages.Present {
		updateFields["messages"] = h.prepareMessages(updateRequest.Messages.Value)
	}
	if updateRequest.Name.Present {
		updateFields["name"] = updateRequest.Name.Value
	}
	if updateRequest.Description.Present {
		updateFields["description"] = updateRequest.Description.Value
	}
	if updateRequest.AgentID.Present {
		updateFields["agent_id"] = updateRequest.AgentID.Value
	}
	if updateRequest.Collaborators.Present {
		updateFields["collaborators"] = updateRequest.Collaborators.Value
	}
	if updateRequest.Action.Present {
		updateFields["action"] = updateRequest.Action.Value
	}
	if updateRequest.Share.Present {
		updateFields["share"] = updateRequest.Share.Value
	}
	if updateRequest.IsArchived.Present {
		updateFields["is_archived"] = updateRequest.IsArchived.Value
	}

	return updateFields
}

func (h *CaseHandler) getObjectIDFromVars(r *http.Request, key string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)[key]))
}

func (h *CaseHandler) getObjectIDFromHeader(r *http.Request, key string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(strings.TrimSpace(r.Header.Get(key)))
}

func (h *CaseHandler) getCaseAndCollaboratorIDs(r *http.Request) (primitive.ObjectID, primitive.ObjectID, error) {
	caseID, err := h.getObjectIDFromVars(r, "id")
	if err != nil {
		return primitive.NilObjectID, primitive.NilObjectID, err
	}
	collaboratorID, err := h.getObjectIDFromVars(r, "collaboratorID")
	if err != nil {
		return primitive.NilObjectID, primitive.NilObjectID, err
	}
	return caseID, collaboratorID, nil
}

func (h *CaseHandler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *CaseHandler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
