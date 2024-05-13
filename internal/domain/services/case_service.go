package services

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CaseService struct {
	caseRepo *repositories.CaseRepository
}

func NewCaseService(caseRepo *repositories.CaseRepository) *CaseService {
	return &CaseService{
		caseRepo: caseRepo,
	}
}

func (s *CaseService) GetAllCases(ctx context.Context) ([]models.Case, error) {
	return s.caseRepo.GetAllCases(ctx)
}

func (s *CaseService) GetCaseByID(ctx context.Context, id string) (models.Case, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Case{}, err
	}
	return s.caseRepo.GetCaseByID(ctx, objectID)
}

func (s *CaseService) GetCasesByCreatorID(ctx context.Context, creatorID string) ([]models.Case, error) {
	return s.caseRepo.GetCasesByCreatorID(ctx, creatorID)
}

func (s *CaseService) CreateCase(ctx context.Context, caseRequest dtos.CreateCaseRequest) (*mongo.InsertOneResult, error) {
	return s.caseRepo.CreateCase(ctx, caseRequest)
}

func (s *CaseService) UpdateCase(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) (*mongo.UpdateResult, error) {
	return s.caseRepo.UpdateCase(ctx, id, updates)
}

func (s *CaseService) DeleteCase(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.caseRepo.DeleteCase(ctx, objectID)
}

func (s *CaseService) AddCollaboratorToCase(ctx context.Context, id, collaboratorID string) (*mongo.UpdateResult, error) {
	caseID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	collaboratorObjectID, err := primitive.ObjectIDFromHex(collaboratorID)
	if err != nil {
		return nil, err
	}
	return s.caseRepo.AddCollaboratorToCase(ctx, caseID, collaboratorObjectID)
}

func (s *CaseService) RemoveCollaboratorFromCase(ctx context.Context, id, collaboratorID string) (*mongo.UpdateResult, error) {
	caseID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	collaboratorObjectID, err := primitive.ObjectIDFromHex(collaboratorID)
	if err != nil {
		return nil, err
	}
	return s.caseRepo.RemoveCollaboratorFromCase(ctx, caseID, collaboratorObjectID)
}
