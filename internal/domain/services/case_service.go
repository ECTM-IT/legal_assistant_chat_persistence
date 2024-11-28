package services

import (
	"context"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
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

// AddDocumentToCase adds a document to a case.
func (s *CaseServiceImpl) AddDocumentToCase(ctx context.Context, caseID primitive.ObjectID, document *models.Document) (*dtos.CaseResponse, error) {
	s.logger.Info("Service Level: Attempting to add document to case")

	_, err := s.caseRepo.AddDocument(ctx, caseID, document)
	if err != nil {
		s.logger.Error("Service Level: Failed to add document to case", err)
		return nil, err
	}

	updatedCase, err := s.GetCaseByID(ctx, caseID)
	if err != nil {
		s.logger.Error("Service Level: Failed to retrieve updated case", err)
		return nil, err
	}

	s.logger.Info("Service Level: Successfully added document to case")
	return updatedCase, nil
}

// UpdateDocument updates a document to a case.
func (s *CaseServiceImpl) UpdateDocument(ctx context.Context, caseID primitive.ObjectID, documentID primitive.ObjectID, document *models.Document) (*dtos.CaseResponse, error) {
	s.logger.Info("Service Level: Attempting to update document to case")

	_, err := s.caseRepo.UpdateDocument(ctx, caseID, documentID, document)
	if err != nil {
		s.logger.Error("Service Level: Failed to update document to case", err)
		return nil, err
	}

	updatedCase, err := s.GetCaseByID(ctx, caseID)
	if err != nil {
		s.logger.Error("Service Level: Failed to retrieve updated case", err)
		return nil, err
	}

	s.logger.Info("Service Level: Successfully updated document to case")
	return updatedCase, nil
}

// UpdateDocument updates a document to a case.
func (s *CaseServiceImpl) AddDocumentCollaborator(ctx context.Context, caseID primitive.ObjectID, documentID primitive.ObjectID, collaborator *models.DocumentCollaborator) (*dtos.CaseResponse, error) {
	s.logger.Info("Service Level: Attempting to update document to case")

	_, err := s.caseRepo.AddDocumentCollaborator(ctx, caseID, documentID, collaborator)
	if err != nil {
		s.logger.Error("Service Level: Failed to update document to case", err)
		return nil, err
	}

	updatedCase, err := s.GetCaseByID(ctx, caseID)
	if err != nil {
		s.logger.Error("Service Level: Failed to retrieve updated case", err)
		return nil, err
	}

	s.logger.Info("Service Level: Successfully updated document to case")
	return updatedCase, nil
}

// DeleteDocumentFromCase removes a document from a case.
func (s *CaseServiceImpl) DeleteDocumentFromCase(ctx context.Context, caseID, documentID primitive.ObjectID) (*dtos.CaseResponse, error) {
	s.logger.Info("Service Level: Attempting to delete document from case")

	_, err := s.caseRepo.DeleteDocument(ctx, caseID, documentID)
	if err != nil {
		s.logger.Error("Service Level: Failed to delete document from case", err)
		return nil, err
	}

	updatedCase, err := s.GetCaseByID(ctx, caseID)
	if err != nil {
		s.logger.Error("Service Level: Failed to retrieve updated case", err)
		return nil, err
	}

	s.logger.Info("Service Level: Successfully deleted document from case")
	return updatedCase, nil
}

// AddFeedbackToMessage adds feedback to a message within a case.
func (s *CaseServiceImpl) AddFeedbackToMessage(ctx context.Context, caseID primitive.ObjectID, messageID string, req *dtos.AddFeedbackRequest) (*models.Feedback, error) {
	s.logger.Info("Service Level: Attempting to add feedback to message in case")

	// Ensure the creator exists (assuming there's a userRepo for this purpose)
	creator, err := s.userRepo.FindUserByID(ctx, req.CreatorID)
	if err != nil {
		s.logger.Error("Service Level: Failed to find feedback creator", err)
		return nil, err
	}

	dateCreated := time.Now()

	// Construct the feedback model from the request
	feedback := models.Feedback{
		ID:           primitive.NewObjectID(),
		CaseID:       caseID,
		MessageID:    messageID,
		CreatorID:    creator.ID,
		Score:        req.Score,
		Reasons:      req.Reasons,
		Comment:      req.Comment,
		CreationDate: dateCreated,
	}

	// Add the feedback to the message in the case
	_, err = s.caseRepo.AddFeedbackToMessage(ctx, caseID, messageID, feedback)
	if err != nil {
		s.logger.Error("Service Level: Failed to add feedback to message in case", err)
		return nil, err
	}

	s.logger.Info("Service Level: Successfully added feedback to message in case")
	return &feedback, nil
}

// GetFeedbackByUserAndMessage retrieves feedback provided by a specific user for a specific message.
func (s *CaseServiceImpl) GetFeedbackByUserAndMessage(ctx context.Context, creatorID primitive.ObjectID, messageID string) ([]models.Feedback, error) {
	s.logger.Info("Service Level: Attempting to retrieve feedback by user and message")

	// Call repository layer to get feedback
	feedbacks, err := s.caseRepo.GetFeedbackByUserAndMessage(ctx, creatorID, messageID)
	if err != nil {
		s.logger.Error("Service Level: Failed to retrieve feedback by user and message", err)
		return nil, err
	}

	s.logger.Info("Service Level: Successfully retrieved feedback by user and message")
	return feedbacks, nil
}
