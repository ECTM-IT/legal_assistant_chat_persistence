package repositories

import (
	"context"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CaseRepository struct {
	caseDAO *daos.CaseDAO
}

func NewCaseRepository(caseDAO *daos.CaseDAO) *CaseRepository {
	return &CaseRepository{
		caseDAO: caseDAO,
	}
}

func (r *CaseRepository) GetAllCases(ctx context.Context) ([]models.Case, error) {
	return r.caseDAO.FindAll(ctx)
}

func (r *CaseRepository) GetCaseByID(ctx context.Context, id primitive.ObjectID) (models.Case, error) {
	return r.caseDAO.FindByID(ctx, id)
}

func (r *CaseRepository) GetCasesByCreatorID(ctx context.Context, creatorID string) ([]models.Case, error) {
	objectID, err := primitive.ObjectIDFromHex(creatorID)
	if err != nil {
		return nil, err
	}
	return r.caseDAO.FindByCreatorID(ctx, objectID)
}

func (r *CaseRepository) CreateCase(ctx context.Context, caseRequest dtos.CreateCaseRequest) (*mongo.InsertOneResult, error) {
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
	return r.caseDAO.Create(ctx, caseModel)
}

func (r *CaseRepository) UpdateCase(ctx context.Context, id primitive.ObjectID, updates dtos.UpdateCaseRequest) (*mongo.UpdateResult, error) {
	return r.caseDAO.Update(ctx, id, updates)
}

func (r *CaseRepository) DeleteCase(ctx context.Context, id primitive.ObjectID) error {
	return r.caseDAO.Delete(ctx, id)
}

func (r *CaseRepository) AddCollaboratorToCase(ctx context.Context, id primitive.ObjectID, collaboratorID primitive.ObjectID) (*mongo.UpdateResult, error) {
	return r.caseDAO.AddCollaborator(ctx, id, collaboratorID)
}

func (r *CaseRepository) RemoveCollaboratorFromCase(ctx context.Context, id primitive.ObjectID, collaboratorID primitive.ObjectID) (*mongo.UpdateResult, error) {
	return r.caseDAO.RemoveCollaborator(ctx, id, collaboratorID)
}
