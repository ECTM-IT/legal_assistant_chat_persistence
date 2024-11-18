package services

import (
	"context"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services/mappers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
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
	AddAgentSkillToCase(ctx context.Context, id primitive.ObjectID, agentSkillRequest dtos.AddAgentSkillToCaseRequest) (dtos.CaseResponse, error)
	RemoveAgentSkillFromCase(ctx context.Context, id, agentSkillID primitive.ObjectID) (dtos.CaseResponse, error)
}

// CaseServiceImpl implements the CaseService interface.
type CaseServiceImpl struct {
	caseRepo   *repositories.CaseRepository
	userRepo   *repositories.UserRepositoryImpl
	mapper     *mappers.CaseConversionServiceImpl
	userMapper *mappers.UserConversionServiceImpl
	logger     logs.Logger
}

// NewCaseService creates a new instance of the case service.
func NewCaseService(caseRepo *repositories.CaseRepository, mapper *mappers.CaseConversionServiceImpl, userMapper *mappers.UserConversionServiceImpl, userRepo *repositories.UserRepositoryImpl, logger logs.Logger) *CaseServiceImpl {
	return &CaseServiceImpl{
		caseRepo:   caseRepo,
		userRepo:   userRepo,
		mapper:     mapper,
		userMapper: userMapper,
		logger:     logger,
	}
}

// GetAllCases retrieves all cases.
func (s *CaseServiceImpl) GetAllCases(ctx context.Context) ([]dtos.CaseResponse, error) {
	s.logger.Info("Service Level: Attempting to retrieve all cases")
	caseModel, err := s.caseRepo.GetAllCases(ctx)
	if err != nil {
		s.logger.Error("Service Level: Failed to retrieve all cases", err)
		return nil, err
	}
	caseResponses := s.mapper.CasesToDTO(caseModel)
	s.logger.Info("Service Level: Successfully retrieved all cases")
	return caseResponses, nil
}

// GetCaseByID retrieves a case by its ID.
func (s *CaseServiceImpl) GetCaseByID(ctx context.Context, id primitive.ObjectID) (*dtos.CaseResponse, error) {
	s.logger.Info("Service Level: Attempting to retrieve case by ID")
	caseModel, err := s.caseRepo.GetCaseByID(ctx, id)
	if err != nil {
		s.logger.Error("Service Level: Failed to retrieve case by ID", err)
		return nil, err
	}
	caseResponse := s.mapper.CaseToDTO(&caseModel)
	s.logger.Info("Service Level: Successfully retrieved case by ID")
	return caseResponse, nil
}

// GetCasesByCreatorID retrieves cases by the creator's ID.
func (s *CaseServiceImpl) GetCasesByCreatorID(ctx context.Context, creatorID primitive.ObjectID) ([]dtos.CaseResponse, error) {
	s.logger.Info("Service Level: Attempting to retrieve cases by creator ID")
	caseModel, err := s.caseRepo.GetCasesByCreatorID(ctx, creatorID)
	if err != nil {
		s.logger.Error("Service Level: Failed to retrieve cases by creator ID", err)
		return nil, err
	}
	caseResponses := s.mapper.CasesToDTO(caseModel)
	s.logger.Info("Service Level: Successfully retrieved cases by creator ID")
	return caseResponses, nil
}

// CreateCase creates a new case.
func (s *CaseServiceImpl) CreateCase(ctx context.Context, caseRequest dtos.CreateCaseRequest) (*dtos.CaseResponse, error) {
	s.logger.Info("Service Level: Attempting to create new case")
	caseModel, err := s.mapper.DTOToCase(caseRequest)
	caseModel.CreationDate = time.Now()
	if err != nil {
		s.logger.Error("Service Level: Failed to convert DTO to case model", err)
		return nil, err
	}
	insertResult, err := s.caseRepo.CreateCase(ctx, *caseModel)
	if err != nil {
		s.logger.Error("Service Level: Failed to create case", err)
		return nil, err
	}
	user, err := s.userRepo.FindUserByID(ctx, caseModel.CreatorID)
	if err != nil {
		s.logger.Error("Service Level: Failed to find user by ID", err)
		return nil, err
	}
	caseIds := make(map[string]interface{})
	caseIds["case_ids"] = append(user.CaseIDs, caseModel.ID)
	_, err = s.userRepo.UpdateUser(ctx, caseModel.CreatorID, caseIds)
	if err != nil {
		s.logger.Error("Service Level: Failed to update user", err)
		return nil, err
	}
	createdCase, err := s.GetCaseByID(ctx, insertResult.InsertedID.(primitive.ObjectID))
	if err != nil {
		s.logger.Error("Service Level: Failed to retrieve created case", err)
		return nil, err
	}
	s.logger.Info("Service Level: Successfully created new case")
	return createdCase, nil
}

// UpdateCase updates an existing case.
func (s *CaseServiceImpl) UpdateCase(ctx context.Context, id primitive.ObjectID, updates dtos.UpdateCaseRequest) (*dtos.CaseResponse, error) {
	s.logger.Info("Service Level: Attempting to update case")
	updateCaseMap, err := s.mapper.UpdateCaseFieldsToMap(updates)
	if err != nil {
		s.logger.Error("Service Level: Failed to map case", err)
		return nil, err
	}
	_, err = s.caseRepo.UpdateCase(ctx, id, updateCaseMap)
	if err != nil {
		s.logger.Error("Service Level: Failed to update case", err)
		return nil, err
	}
	updatedCase, err := s.GetCaseByID(ctx, id)
	if err != nil {
		s.logger.Error("Service Level: Failed to retrieve updated case", err)
		return nil, err
	}
	s.logger.Info("Service Level: Successfully updated case")
	return updatedCase, nil
}

// DeleteCase deletes a case by its ID.
func (s *CaseServiceImpl) DeleteCase(ctx context.Context, id primitive.ObjectID) (*dtos.CaseResponse, error) {
	s.logger.Info("Service Level: Attempting to delete case")
	deletedCase, err := s.GetCaseByID(ctx, id)
	if err != nil {
		s.logger.Error("Service Level: Failed to retrieve case for deletion", err)
		return nil, err
	}
	err = s.caseRepo.DeleteCase(ctx, id)
	if err != nil {
		s.logger.Error("Service Level: Failed to delete case", err)
		return nil, err
	}
	s.logger.Info("Service Level: Successfully deleted case")
	return deletedCase, nil
}

// AddCollaboratorToCase adds a collaborator to a case.
func (s *CaseServiceImpl) AddCollaboratorToCase(ctx context.Context, id primitive.ObjectID, email string, canEdit bool) (*dtos.UserResponse, error) {
	s.logger.Info("Service Level: Attempting to add collaborator to case")
	collaborator, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error("Service Level: Failed to find collaborator by email", err)
		return nil, err
	}

	updates := map[string]interface{}{
		"collaborator_id": collaborator.ID,
		"can_edit":        canEdit,
	}
	_, err = s.caseRepo.AddCollaboratorToCase(ctx, id, updates)
	if err != nil {
		s.logger.Error("Service Level: Failed to add collaborator to case", err)
		return nil, err
	}
	s.logger.Info("Service Level: Successfully added collaborator to case")
	return s.userMapper.UserToDTO(collaborator), nil
}

// RemoveCollaboratorFromCase removes a collaborator from a case.
func (s *CaseServiceImpl) RemoveCollaboratorFromCase(ctx context.Context, id, collaboratorID primitive.ObjectID) (*dtos.CaseResponse, error) {
	s.logger.Info("Service Level: Attempting to remove collaborator from case")
	_, err := s.caseRepo.RemoveCollaboratorFromCase(ctx, id, collaboratorID)
	if err != nil {
		s.logger.Error("Service Level: Failed to remove collaborator from case", err)
		return nil, err
	}
	updatedCase, err := s.GetCaseByID(ctx, id)
	if err != nil {
		s.logger.Error("Service Level: Failed to retrieve updated case", err)
		return nil, err
	}
	s.logger.Info("Service Level: Successfully removed collaborator from case")
	return updatedCase, nil
}

// AddAgentSkillToCase adds a agent skill to a case.
func (s *CaseServiceImpl) AddAgentSkillToCase(ctx context.Context, id primitive.ObjectID, agentSkillRequest dtos.AddAgentSkillToCaseRequest) (*dtos.CaseResponse, error) {
	s.logger.Info("Service Level: Attempting to add agent skill to case")

	updates := map[string]interface{}{
		"id":       agentSkillRequest.ID,
		"agent_id": agentSkillRequest.AgentID,
		"name":     agentSkillRequest.Name,
	}
	_, err := s.caseRepo.AddAgentSkillToCase(ctx, id, updates)
	if err != nil {
		s.logger.Error("Service Level: Failed to add agent skill to case", err)
		return nil, err
	}
	updatedCase, err := s.GetCaseByID(ctx, id)
	if err != nil {
		s.logger.Error("Service Level: Failed to retrieve updated case", err)
		return nil, err
	}
	s.logger.Info("Service Level: Successfully added agent skill to case")
	return updatedCase, nil
}

// RemoveAgentSkillFromCase removes a agent skill from a case.
func (s *CaseServiceImpl) RemoveAgentSkillFromCase(ctx context.Context, id, agentSkillID primitive.ObjectID) (*dtos.CaseResponse, error) {
	s.logger.Info("Service Level: Attempting to remove agent skill from case")
	_, err := s.caseRepo.RemoveAgentSkillFromCase(ctx, id, agentSkillID)
	if err != nil {
		s.logger.Error("Service Level: Failed to remove agent skill from case", err)
		return nil, err
	}
	updatedCase, err := s.GetCaseByID(ctx, id)
	if err != nil {
		s.logger.Error("Service Level: Failed to retrieve updated case", err)
		return nil, err
	}
	s.logger.Info("Service Level: Successfully removed agent skill from case")
	return updatedCase, nil
}
