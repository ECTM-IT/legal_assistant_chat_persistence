package repositories

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
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

func (r *CaseRepository) GetCasesByCreatorID(ctx context.Context, creatorID primitive.ObjectID) ([]models.Case, error) {
	return r.caseDAO.FindByCreatorID(ctx, creatorID)
}

func (r *CaseRepository) CreateCase(ctx context.Context, caseModel models.Case) (*mongo.InsertOneResult, error) {
	return r.caseDAO.Create(ctx, &caseModel)
}

func (r *CaseRepository) UpdateCase(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) (*mongo.UpdateResult, error) {
	return r.caseDAO.Update(ctx, id, updates)
}

func (r *CaseRepository) DeleteCase(ctx context.Context, id primitive.ObjectID) error {
	return r.caseDAO.Delete(ctx, id)
}

func (r *CaseRepository) AddCollaboratorToCase(ctx context.Context, id primitive.ObjectID, collaborator map[string]interface{}) (*mongo.UpdateResult, error) {
	return r.caseDAO.AddCollaborator(ctx, id, collaborator)
}

func (r *CaseRepository) RemoveCollaboratorFromCase(ctx context.Context, id primitive.ObjectID, collaboratorID primitive.ObjectID) (*mongo.UpdateResult, error) {
	return r.caseDAO.RemoveCollaborator(ctx, id, collaboratorID)
}

func (r *CaseRepository) AddDocument(ctx context.Context, caseID primitive.ObjectID, document *models.Document) (*mongo.UpdateResult, error) {
	return r.caseDAO.AddDocument(ctx, caseID, document)
}

func (r *CaseRepository) DeleteDocument(ctx context.Context, caseID, documentID primitive.ObjectID) (*mongo.UpdateResult, error) {
	return r.caseDAO.DeleteDocument(ctx, caseID, documentID)
}

func (r *CaseRepository) AddFeedbackToMessage(ctx context.Context, feedback models.Feedback) (*mongo.UpdateResult, error) {
	return r.caseDAO.AddFeedback(ctx, feedback.CaseID, feedback.MessageID, feedback)
}

func (r *CaseRepository) GetFeedbackByUserAndMessage(ctx context.Context, creatorID, messageID primitive.ObjectID) ([]models.Feedback, error) {
	return r.caseDAO.GetFeedbackByUserAndMessage(ctx, creatorID, messageID)
}
