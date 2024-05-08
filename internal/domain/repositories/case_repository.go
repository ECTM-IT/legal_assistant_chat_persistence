package repositories

import (
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CaseRepository struct {
	caseDAO *daos.CaseDAO
}

func NewCaseRepository(caseDAO *daos.CaseDAO) *CaseRepository {
	return &CaseRepository{
		caseDAO: caseDAO,
	}
}

func (r *CaseRepository) GetAllCases() ([]dtos.CaseResponse, error) {
	return r.caseDAO.FindAll()
}

func (r *CaseRepository) GetCaseByID(id primitive.ObjectID) (dtos.CaseResponse, error) {
	return r.caseDAO.FindByID(id)
}

func (r *CaseRepository) GetCasesByCreatorID(creatorID string) ([]dtos.CaseResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(creatorID)
	if err != nil {
		return nil, err
	}
	return r.caseDAO.FindByCreatorID(objectID)
}

func (r *CaseRepository) CreateCase(caseRequest dtos.CreateCaseRequest) error {
	messageDto := caseRequest.Messages
	messages := []models.Message{}
	if messageDto.Present {
		for _, msg := range messageDto.Val {
			messageModel := models.Message{
				Content:     msg.Content.OrElse(""),
				SenderID:    msg.SenderID.OrElse(primitive.NilObjectID),
				RecipientID: msg.RecipientID.OrElse(primitive.NilObjectID),
				Skill:       msg.Skill.OrElse(""),
			}
			messages = append(messages, messageModel)
		}
	}
	caseModel := &models.Case{
		ID:              primitive.NewObjectID(),
		Name:            caseRequest.Name.OrElse(""),
		CreatorID:       caseRequest.CreatorID.OrElse(primitive.NilObjectID),
		Messages:        messages,
		CollaboratorIDs: caseRequest.CollaboratorIDs.OrElse([]primitive.ObjectID{}),
		Action:          caseRequest.Action.OrElse(""),
		AgentID:         caseRequest.AgentID.OrElse(primitive.NilObjectID),
		LastEdit:        caseRequest.LastEdit.OrElse(time.Now()),
		Share:           caseRequest.Share.OrElse(false),
		IsArchived:      caseRequest.IsArchived.OrElse(false),
	}
	return r.caseDAO.Create(caseModel)
}

func (r *CaseRepository) UpdateCase(id primitive.ObjectID, updates dtos.UpdateCaseRequest) error {
	return r.caseDAO.Update(id, updates)
}

func (r *CaseRepository) DeleteCase(id primitive.ObjectID) error {
	return r.caseDAO.Delete(id)
}

func (r *CaseRepository) AddCollaboratorToCase(id primitive.ObjectID, collaboratorID primitive.ObjectID) error {
	return r.caseDAO.AddCollaborator(id, collaboratorID)
}

func (r *CaseRepository) RemoveCollaboratorFromCase(id primitive.ObjectID, collaboratorID primitive.ObjectID) error {
	return r.caseDAO.RemoveCollaborator(id, collaboratorID)
}
