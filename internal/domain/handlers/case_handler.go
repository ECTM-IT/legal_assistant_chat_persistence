package handlers

import (
	"context"
	"net/http"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
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
	AddDocumentToCase(ctx context.Context, caseID primitive.ObjectID, document *models.Document) (*dtos.CaseResponse, error)
	UpdateDocument(ctx context.Context, caseID primitive.ObjectID, documentID primitive.ObjectID, document *models.Document) (*dtos.CaseResponse, error)
	DeleteDocumentFromCase(ctx context.Context, caseID, documentID primitive.ObjectID) (*dtos.CaseResponse, error)
	AddAgentSkillToCase(ctx context.Context, caseID primitive.ObjectID, agentSkill *dtos.AddAgentSkillToCaseRequest) (*dtos.CaseResponse, error)
	DeleteAgentSkillFromCase(ctx context.Context, caseID, agentSkillID primitive.ObjectID) (*dtos.CaseResponse, error)
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

	var req dtos.AddCollaboratorToCase
	if err := h.DecodeJSONBody(r, &req); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	newUser, err := h.service.AddCollaboratorToCase(r.Context(), caseID, req.Email.Value, req.Edit.Value)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to add collaborator to case")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, newUser)
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

func (h *CaseHandler) AddDocumentToCase(w http.ResponseWriter, r *http.Request) {
	caseID, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid case ID")
		return
	}

	var doc dtos.AddDocumentToCase // Assuming a DTO to parse the incoming JSON payload
	if err := h.DecodeJSONBody(r, &doc); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	document := &models.Document{
		Sender:      doc.Sender.OrElse(""),
		FileName:    doc.FileName.OrElse(""),
		FileType:    doc.FileType.OrElse(""),
		FileContent: doc.FileContent.OrElse(""),
	}

	updatedCase, err := h.service.AddDocumentToCase(r.Context(), caseID, document)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to add document to case")
		return
	}

	h.RespondWithJSON(w, http.StatusOK, updatedCase)
}

func (h *CaseHandler) UpdateDocument(w http.ResponseWriter, r *http.Request) {
	caseID, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid case ID")
		return
	}

	documentID, err := h.ParseObjectID(r, "documentID", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	var doc dtos.UpdateDocument // Assuming a DTO to parse the incoming JSON payload
	if err := h.DecodeJSONBody(r, &doc); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	document := &models.Document{
		FileName:    doc.FileName.OrElse(""),
		FileType:    doc.FileType.OrElse(""),
		FileContent: doc.FileContent.OrElse(""),
	}

	updatedCase, err := h.service.UpdateDocument(r.Context(), caseID, documentID, document)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to add document to case")
		return
	}

	h.RespondWithJSON(w, http.StatusOK, updatedCase)
}

// AddDocumentCollaborator handles adding a collaborator to a specific document in a case
func (h *CaseHandler) AddDocumentCollaborator(w http.ResponseWriter, r *http.Request) {
	caseID, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid case ID")
		return
	}

	documentID, err := h.ParseObjectID(r, "documentID", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	var collaborator dtos.DocumentCollaboratorRequest
	if err := h.DecodeJSONBody(r, &collaborator); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Create a new collaborator model
	newCollaborator := &models.DocumentCollaborator{
		Email: collaborator.Email.OrElse(""),
		Edit:  collaborator.Edit.OrElse(false),
	}

	// Call the service layer to add the collaborator
	updatedCase, err := h.service.AddDocumentCollaborator(r.Context(), caseID, documentID, newCollaborator)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to add collaborator to document")
		return
	}

	h.RespondWithJSON(w, http.StatusOK, updatedCase)
}

func (h *CaseHandler) DeleteDocumentFromCase(w http.ResponseWriter, r *http.Request) {
	caseID, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid case ID")
		return
	}

	documentID, err := h.ParseObjectID(r, "documentID", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid document ID")
		return
	}

	updatedCase, err := h.service.DeleteDocumentFromCase(r.Context(), caseID, documentID)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to delete document from case")
		return
	}

	h.RespondWithJSON(w, http.StatusOK, updatedCase)
}

func (h *CaseHandler) AddFeedbackToMessage(w http.ResponseWriter, r *http.Request) {
	// Parse caseID from the URL
	caseID, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid case ID")
		return
	}

	// Parse messageID from the URL
	messageID, err := h.ParseObjectID(r, "messageId", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	// Decode the request body into AddFeedbackRequest DTO
	var req dtos.AddFeedbackRequest
	if err := h.DecodeJSONBody(r, &req); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Set caseID and messageID in request to ensure they are passed correctly
	req.CaseID = caseID
	req.MessageID = messageID.String()

	// Call the service layer to add feedback to the message
	feedback, err := h.service.AddFeedbackToMessage(r.Context(), &req)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to add feedback to message")
		return
	}

	// Respond with the created feedback
	h.RespondWithJSON(w, http.StatusCreated, feedback)
}

// GetFeedbackByUserAndMessageHandler handles the HTTP request to retrieve feedback for a specific user and message.
func (h *CaseHandler) GetFeedbackByUserAndMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Parse creatorID from the URL or query parameters
	_, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid case ID")
		return
	}

	// Parse creatorID from the URL or query parameters
	creatorID, err := h.ParseObjectID(r, "creatorId", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid creator ID")
		return
	}

	// Parse messageID from the URL or query parameters
	messageID, err := h.ParseObjectID(r, "messageId", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	// Call the service layer to get the feedback
	feedbacks, err := h.service.GetFeedbackByUserAndMessage(r.Context(), creatorID, messageID.String())
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve feedback")
		return
	}

	// Respond with the feedback
	h.RespondWithJSON(w, http.StatusOK, feedbacks)
}

func (h *CaseHandler) AddAgentSkillToCase(w http.ResponseWriter, r *http.Request) {
	caseID, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid case ID")
		return
	}
	var req dtos.AddAgentSkillToCaseRequest
	if err := h.DecodeJSONBody(r, &req); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	newUser, err := h.service.AddAgentSkillToCase(r.Context(), caseID, req)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to add agent skill to case")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, newUser)
}
func (h *CaseHandler) RemoveAgentSkillFromCase(w http.ResponseWriter, r *http.Request) {
	caseID, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid case ID")
		return
	}
	agentSkillID, err := h.ParseObjectID(r, "agentSkillID", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid agent skill ID")
		return
	}
	updatedCase, err := h.service.RemoveAgentSkillFromCase(r.Context(), caseID, agentSkillID)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to remove agent skill from case")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, updatedCase)
}
