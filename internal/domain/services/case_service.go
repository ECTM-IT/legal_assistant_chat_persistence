package services

import (
	"context"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services/mappers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CaseService defines the operations available for managing cases.
type CaseService interface {
	GetAllCases(ctx context.Context) ([]dtos.CaseResponse, error)
	GetCaseByID(ctx context.Context, id primitive.ObjectID) (dtos.CaseResponse, error)
	GetCasesByCreatorID(ctx context.Context, creatorID primitive.ObjectID) ([]dtos.CaseResponse, error)
	CreateCase(ctx context.Context, caseRequest dtos.CreateCaseRequest) (dtos.CaseResponse, error)
	UpdateCase(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) (dtos.CaseResponse, error)
	DeleteCase(ctx context.Context, id primitive.ObjectID) (dtos.CaseResponse, error)
	AddCollaboratorToCase(ctx context.Context, id, collaboratorID primitive.ObjectID) (dtos.CaseResponse, error)
	RemoveCollaboratorFromCase(ctx context.Context, id, collaboratorID primitive.ObjectID) (dtos.CaseResponse, error)
}

// CaseServiceImpl implements the CaseService interface.
type CaseServiceImpl struct {
	caseRepo *repositories.CaseRepository
	userRepo *repositories.UserRepositoryImpl
	mapper   *mappers.CaseConversionServiceImpl
}

// NewCaseService creates a new instance of the case service.
func NewCaseService(caseRepo *repositories.CaseRepository, mapper *mappers.CaseConversionServiceImpl, userRepo *repositories.UserRepositoryImpl) *CaseServiceImpl {
	return &CaseServiceImpl{
		caseRepo: caseRepo,
		userRepo: userRepo,
		mapper:   mapper,
	}
}

// GetAllCases retrieves all cases.
func (s *CaseServiceImpl) GetAllCases(ctx context.Context) ([]dtos.CaseResponse, error) {
	caseModel, err := s.caseRepo.GetAllCases(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.CasesToDTO(caseModel), nil
}

// GetCaseByID retrieves a case by its ID.
func (s *CaseServiceImpl) GetCaseByID(ctx context.Context, id primitive.ObjectID) (*dtos.CaseResponse, error) {
	caseModel, err := s.caseRepo.GetCaseByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapper.CaseToDTO(&caseModel), nil
}

// GetCasesByCreatorID retrieves cases by the creator's ID.
func (s *CaseServiceImpl) GetCasesByCreatorID(ctx context.Context, creatorID primitive.ObjectID) ([]dtos.CaseResponse, error) {
	caseModel, err := s.caseRepo.GetCasesByCreatorID(ctx, creatorID)
	if err != nil {
		return nil, err
	}
	return s.mapper.CasesToDTO(caseModel), nil
}

// CreateCase creates a new case.
func (s *CaseServiceImpl) CreateCase(ctx context.Context, caseRequest dtos.CreateCaseRequest) (*dtos.CaseResponse, error) {
	caseModel, err := s.mapper.DTOToCase(caseRequest)
	caseModel.CreationDate = time.Now()
	if err != nil {
		return nil, err
	}
	insertResult, err := s.caseRepo.CreateCase(ctx, *caseModel)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.FindUserByID(ctx, caseModel.CreatorID)
	if err != nil {
		return nil, err
	}
	caseIds := make(map[string]interface{})
	caseIds["case_ids"] = append(user.CaseIDs, caseModel.ID)
	_, err = s.userRepo.UpdateUser(ctx, caseModel.CreatorID, caseIds)
	if err != nil {
		return nil, err
	}
	return s.GetCaseByID(ctx, insertResult.InsertedID.(primitive.ObjectID))
}

// UpdateCase updates an existing case.
func (s *CaseServiceImpl) UpdateCase(ctx context.Context, id primitive.ObjectID, updates dtos.UpdateCaseRequest) (*dtos.CaseResponse, error) {
	updateCaseMap := s.mapper.UpdateCaseFieldsToMap(updates)
	_, err := s.caseRepo.UpdateCase(ctx, id, updateCaseMap)
	if err != nil {
		return nil, err
	}
	return s.GetCaseByID(ctx, id)
}

// DeleteCase deletes a case by its ID.
func (s *CaseServiceImpl) DeleteCase(ctx context.Context, id primitive.ObjectID) (*dtos.CaseResponse, error) {
	deletedCase, err := s.GetCaseByID(ctx, id)
	if err != nil {
		return nil, err
	}
	err = s.caseRepo.DeleteCase(ctx, id)
	if err != nil {
		return nil, err
	}
	return deletedCase, nil
}

// AddCollaboratorToCase adds a collaborator to a case.
func (s *CaseServiceImpl) AddCollaboratorToCase(ctx context.Context, id primitive.ObjectID, email string, canEdit bool) (*dtos.CaseResponse, error) {
	collaborator, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"collaborator_id": collaborator.ID,
		"can_edit":        canEdit,
	}
	_, err = s.caseRepo.AddCollaboratorToCase(ctx, id, updates)
	if err != nil {
		return nil, err
	}
	return s.GetCaseByID(ctx, id)
}

// RemoveCollaboratorFromCase removes a collaborator from a case.
func (s *CaseServiceImpl) RemoveCollaboratorFromCase(ctx context.Context, id, collaboratorID primitive.ObjectID) (*dtos.CaseResponse, error) {
	s.caseRepo.RemoveCollaboratorFromCase(ctx, id, collaboratorID)
	return s.GetCaseByID(ctx, id)
}
