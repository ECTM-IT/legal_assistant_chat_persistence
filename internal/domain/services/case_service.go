package services

import (
	"context"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CaseService defines the operations available for managing cases.
type CaseService interface {
	GetAllCases(ctx context.Context) ([]models.Case, error)
	GetCaseByID(ctx context.Context, id primitive.ObjectID) (models.Case, error)
	GetCasesByCreatorID(ctx context.Context, creatorID primitive.ObjectID) ([]models.Case, error)
	CreateCase(ctx context.Context, caseRequest dtos.CreateCaseRequest) (*mongo.InsertOneResult, error)
	UpdateCase(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) (*mongo.UpdateResult, error)
	DeleteCase(ctx context.Context, id primitive.ObjectID) error
	AddCollaboratorToCase(ctx context.Context, id, collaboratorID primitive.ObjectID) (*mongo.UpdateResult, error)
	RemoveCollaboratorFromCase(ctx context.Context, id, collaboratorID primitive.ObjectID) (*mongo.UpdateResult, error)
}

// caseServiceImpl implements the CaseService interface.
type CaseServiceImpl struct {
	caseRepo *repositories.CaseRepository
}

// NewCaseService creates a new instance of the case service.
func NewCaseService(caseRepo *repositories.CaseRepository) *CaseServiceImpl {
	return &CaseServiceImpl{
		caseRepo: caseRepo,
	}
}

// GetAllCases retrieves all cases.
func (s *CaseServiceImpl) GetAllCases(ctx context.Context) ([]models.Case, error) {
	return s.caseRepo.GetAllCases(ctx)
}

// GetCaseByID retrieves a case by its ID.
func (s *CaseServiceImpl) GetCaseByID(ctx context.Context, id primitive.ObjectID) (models.Case, error) {
	return s.caseRepo.GetCaseByID(ctx, id)
}

// GetCasesByCreatorID retrieves cases by the creator's ID.
func (s *CaseServiceImpl) GetCasesByCreatorID(ctx context.Context, creatorID primitive.ObjectID) ([]models.Case, error) {
	return s.caseRepo.GetCasesByCreatorID(ctx, creatorID)
}

// CreateCase creates a new case.
func (s *CaseServiceImpl) CreateCase(ctx context.Context, caseRequest dtos.CreateCaseRequest) (*mongo.InsertOneResult, error) {
	return s.caseRepo.CreateCase(ctx, caseRequest)
}

// UpdateCase updates an existing case.
func (s *CaseServiceImpl) UpdateCase(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) (*mongo.UpdateResult, error) {
	return s.caseRepo.UpdateCase(ctx, id, updates)
}

// DeleteCase deletes a case by its ID.
func (s *CaseServiceImpl) DeleteCase(ctx context.Context, id primitive.ObjectID) error {
	return s.caseRepo.DeleteCase(ctx, id)
}

// AddCollaboratorToCase adds a collaborator to a case.
func (s *CaseServiceImpl) AddCollaboratorToCase(ctx context.Context, id, collaboratorID primitive.ObjectID) (*mongo.UpdateResult, error) {
	return s.caseRepo.AddCollaboratorToCase(ctx, id, collaboratorID)
}

// RemoveCollaboratorFromCase removes a collaborator from a case.
func (s *CaseServiceImpl) RemoveCollaboratorFromCase(ctx context.Context, id, collaboratorID primitive.ObjectID) (*mongo.UpdateResult, error) {
	return s.caseRepo.RemoveCollaboratorFromCase(ctx, id, collaboratorID)
}
