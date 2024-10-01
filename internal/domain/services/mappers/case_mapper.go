package mappers

import (
	"errors"
	"fmt"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CaseConversionService interface {
	DTOToCase(caseRequest dtos.CreateCaseRequest) (*models.Case, error)
	CaseToDTO(caseModel *models.Case) *dtos.CaseResponse
	CasesToDTO(cases []models.Case) []dtos.CaseResponse
	UpdateCaseFieldsToMap(updateRequest dtos.UpdateCaseRequest) map[string]interface{}
	CollaboratorsToDTO(collaborators []models.Collaborators) []dtos.CollaboratorResponse
	DTOToCollaborators(collaboratorsDTO []dtos.CollaboratorResponse) ([]models.Collaborators, error)
	MessageToDTO(message models.Message) dtos.MessageResponse
	DTOToMessage(messageDTO dtos.MessageResponse) (models.Message, error)
	MessagesToDTO(messages []models.Message) []dtos.MessageResponse
	DTOToMessages(messagesDTO []dtos.MessageResponse) ([]models.Message, error)
}

type CaseConversionServiceImpl struct {
	logger logs.Logger
}

func NewCaseConversionService(logger logs.Logger) *CaseConversionServiceImpl {
	return &CaseConversionServiceImpl{
		logger: logger,
	}
}

func (s *CaseConversionServiceImpl) DTOToCase(caseRequest dtos.CreateCaseRequest) (*models.Case, error) {
	s.logger.Info("Converting DTO to Case")
	if !caseRequest.CreatorID.Present {
		s.logger.Error("Failed to convert DTO to Case: creator ID is required", errors.New("creator ID is required"))
		return nil, errors.New("creator ID is required")
	}

	messages, err := s.DTOToMessages(caseRequest.Messages.OrElse(nil))
	if err != nil {
		s.logger.Error("Failed to convert messages", err)
		return nil, fmt.Errorf("error converting messages: %w", err)
	}

	collaborators, err := s.DTOToCollaborators(caseRequest.Collaborators.OrElse(nil))
	if err != nil {
		s.logger.Error("Failed to convert collaborators", err)
		return nil, fmt.Errorf("error converting collaborators: %w", err)
	}

	now := time.Now()
	caseModel := &models.Case{
		ID:            primitive.NewObjectID(),
		Name:          caseRequest.Name.OrElse("New Case"),
		CreatorID:     caseRequest.CreatorID.Value,
		Messages:      messages,
		Collaborators: collaborators,
		Action:        caseRequest.Action.OrElse("summarize"),
		AgentID:       caseRequest.AgentID.OrElse(primitive.NilObjectID),
		LastEdit:      caseRequest.LastEdit.OrElse(now),
		CreationDate:  now,
		Share:         caseRequest.Share.OrElse(false),
		IsArchived:    caseRequest.IsArchived.OrElse(false),
	}
	s.logger.Info("Successfully converted DTO to Case")
	return caseModel, nil
}

func (s *CaseConversionServiceImpl) CaseToDTO(caseModel *models.Case) *dtos.CaseResponse {
	s.logger.Info("Converting Case to DTO")
	if caseModel == nil {
		s.logger.Warn("Attempted to convert nil Case to DTO")
		return nil
	}

	dto := &dtos.CaseResponse{
		ID:            helpers.NewNullable(caseModel.ID),
		Name:          helpers.NewNullable(caseModel.Name),
		CreatorID:     helpers.NewNullable(caseModel.CreatorID),
		Messages:      helpers.NewNullable(s.MessagesToDTO(caseModel.Messages)),
		Collaborators: helpers.NewNullable(s.CollaboratorsToDTO(caseModel.Collaborators)),
		Action:        helpers.NewNullable(caseModel.Action),
		AgentID:       helpers.NewNullable(caseModel.AgentID),
		LastEdit:      helpers.NewNullable(caseModel.LastEdit),
		CreationDate:  helpers.NewNullable(caseModel.CreationDate),
		Share:         helpers.NewNullable(caseModel.Share),
		IsArchived:    helpers.NewNullable(caseModel.IsArchived),
	}
	s.logger.Info("Successfully converted Case to DTO")
	return dto
}

func (s *CaseConversionServiceImpl) CasesToDTO(cases []models.Case) []dtos.CaseResponse {
	s.logger.Info("Converting multiple Cases to DTOs")
	if cases == nil {
		s.logger.Warn("Attempted to convert nil Cases slice to DTOs")
		return nil
	}
	caseDTOs := make([]dtos.CaseResponse, 0, len(cases))
	for _, caseModel := range cases {
		if dto := s.CaseToDTO(&caseModel); dto != nil {
			caseDTOs = append(caseDTOs, *dto)
		}
	}
	s.logger.Info("Successfully converted multiple Cases to DTOs")
	return caseDTOs
}

func (s *CaseConversionServiceImpl) UpdateCaseFieldsToMap(updateRequest dtos.UpdateCaseRequest) map[string]interface{} {
	s.logger.Info("Converting UpdateCaseRequest to map")
	updateFields := make(map[string]interface{})

	if updateRequest.Name.Present {
		updateFields["name"] = updateRequest.Name.Value
	}
	if updateRequest.Messages.Present {
		if messages, err := s.DTOToMessages(updateRequest.Messages.Value); err == nil {
			updateFields["messages"] = messages
		} else {
			s.logger.Error("Failed to convert messages for update", err)
		}
	}
	if updateRequest.Collaborators.Present {
		if collaborators, err := s.DTOToCollaborators(updateRequest.Collaborators.Value); err == nil {
			updateFields["collaborators"] = collaborators
		} else {
			s.logger.Error("Failed to convert collaborators for update", err)
		}
	}
	if updateRequest.Action.Present {
		updateFields["action"] = updateRequest.Action.Value
	}
	if updateRequest.AgentID.Present {
		updateFields["agent_id"] = updateRequest.AgentID.Value
	}
	if updateRequest.Share.Present {
		updateFields["share"] = updateRequest.Share.Value
	}
	if updateRequest.IsArchived.Present {
		updateFields["is_archived"] = updateRequest.IsArchived.Value
	}

	updateFields["last_edit"] = time.Now()

	s.logger.Info("Successfully converted UpdateCaseRequest to map")
	return updateFields
}

func (s *CaseConversionServiceImpl) CollaboratorsToDTO(collaborators []models.Collaborators) []dtos.CollaboratorResponse {
	s.logger.Info("Converting Collaborators to DTOs")
	if collaborators == nil {
		s.logger.Warn("Attempted to convert nil Collaborators slice to DTOs")
		return nil
	}

	collaboratorDTOs := make([]dtos.CollaboratorResponse, 0, len(collaborators))
	for _, collaborator := range collaborators {
		collaboratorDTOs = append(collaboratorDTOs, dtos.CollaboratorResponse{
			ID:   helpers.NewNullable(collaborator.ID),
			Edit: helpers.NewNullable(collaborator.Edit),
		})
	}
	s.logger.Info("Successfully converted Collaborators to DTOs")
	return collaboratorDTOs
}

func (s *CaseConversionServiceImpl) DTOToCollaborators(collaboratorsDTO []dtos.CollaboratorResponse) ([]models.Collaborators, error) {
	s.logger.Info("Converting DTOs to Collaborators")
	if collaboratorsDTO == nil {
		s.logger.Warn("Attempted to convert nil CollaboratorResponse slice to Collaborators")
		return nil, nil
	}

	collaborators := make([]models.Collaborators, 0, len(collaboratorsDTO))
	for _, dto := range collaboratorsDTO {
		if !dto.ID.Present {
			s.logger.Error("Failed to convert DTO to Collaborator: collaborator ID is required", errors.New("collaborator ID is required"))
			return nil, errors.New("collaborator ID is required")
		}
		collaborators = append(collaborators, models.Collaborators{
			ID:   dto.ID.Value,
			Edit: dto.Edit.OrElse(false),
		})
	}
	s.logger.Info("Successfully converted DTOs to Collaborators")
	return collaborators, nil
}

func (s *CaseConversionServiceImpl) MessageToDTO(message models.Message) dtos.MessageResponse {
	s.logger.Info("Converting Message to DTO")
	dto := dtos.MessageResponse{
		Content:      helpers.NewNullable(message.Content),
		Sender:       helpers.NewNullable(message.Sender),
		Recipient:    helpers.NewNullable(message.Recipient),
		FunctionCall: helpers.NewNullable(message.FunctionCall),
		DocumentPath: helpers.NewNullable(message.DocumentPath),
	}
	s.logger.Info("Successfully converted Message to DTO")
	return dto
}

func (s *CaseConversionServiceImpl) DTOToMessage(messageDTO dtos.MessageResponse) (models.Message, error) {
	s.logger.Info("Converting DTO to Message")
	if !messageDTO.Content.Present || !messageDTO.Sender.Present || !messageDTO.Recipient.Present {
		s.logger.Error("Failed to convert DTO to Message: content, sender, and recipient are required", errors.New("content, sender, and recipient are required"))
		return models.Message{}, errors.New("content, sender, and recipient are required")
	}
	message := models.Message{
		Content:      messageDTO.Content.Value,
		Sender:       messageDTO.Sender.Value,
		Recipient:    messageDTO.Recipient.Value,
		FunctionCall: messageDTO.FunctionCall.OrElse(false),
		DocumentPath: messageDTO.DocumentPath.OrElse(""),
	}
	s.logger.Info("Successfully converted DTO to Message")
	return message, nil
}

func (s *CaseConversionServiceImpl) MessagesToDTO(messages []models.Message) []dtos.MessageResponse {
	s.logger.Info("Converting multiple Messages to DTOs")
	if messages == nil {
		s.logger.Warn("Attempted to convert nil Messages slice to DTOs")
		return nil
	}
	messageDTOs := make([]dtos.MessageResponse, len(messages))
	for i, message := range messages {
		messageDTOs[i] = s.MessageToDTO(message)
	}
	s.logger.Info("Successfully converted multiple Messages to DTOs")
	return messageDTOs
}

func (s *CaseConversionServiceImpl) DTOToMessages(messagesDTO []dtos.MessageResponse) ([]models.Message, error) {
	s.logger.Info("Converting DTOs to Messages")
	if messagesDTO == nil {
		s.logger.Warn("Attempted to convert nil MessageResponse slice to Messages")
		return nil, nil
	}
	messages := make([]models.Message, 0, len(messagesDTO))
	for _, dto := range messagesDTO {
		message, err := s.DTOToMessage(dto)
		if err != nil {
			s.logger.Error("Failed to convert DTO to Message", err)
			return nil, fmt.Errorf("error converting message: %w", err)
		}
		messages = append(messages, message)
	}
	s.logger.Info("Successfully converted DTOs to Messages")
	return messages, nil
}
