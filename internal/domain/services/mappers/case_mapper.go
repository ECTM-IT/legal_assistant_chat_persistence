package mappers

import (
	"errors"
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
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

type CaseConversionServiceImpl struct{}

func NewCaseConversionService() *CaseConversionServiceImpl {
	return &CaseConversionServiceImpl{}
}

func (s *CaseConversionServiceImpl) DTOToCase(caseRequest dtos.CreateCaseRequest) (*models.Case, error) {
	if !caseRequest.CreatorID.Present {
		return nil, errors.New("creator ID is required")
	}

	messages, err := s.DTOToMessages(caseRequest.Messages.Value)
	if err != nil {
		return nil, err
	}

	collaborators, err := s.DTOToCollaborators(caseRequest.Collaborators.Value)
	if err != nil {
		return nil, err
	}

	return &models.Case{
		ID:            primitive.NewObjectID(),
		Name:          caseRequest.Name.OrElse("New Case"),
		CreatorID:     caseRequest.CreatorID.Value,
		Messages:      messages,
		Collaborators: collaborators,
		Action:        caseRequest.Action.OrElse("summarize"),
		AgentID:       caseRequest.AgentID.OrElse(primitive.NilObjectID),
		LastEdit:      caseRequest.LastEdit.OrElse(time.Now()),
		CreationDate:  time.Now(),
		Share:         caseRequest.Share.OrElse(false),
		IsArchived:    caseRequest.IsArchived.OrElse(false),
	}, nil
}

func (s *CaseConversionServiceImpl) CaseToDTO(caseModel *models.Case) *dtos.CaseResponse {
	if caseModel == nil {
		return nil
	}

	return &dtos.CaseResponse{
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
}

func (s *CaseConversionServiceImpl) CasesToDTO(cases []models.Case) []dtos.CaseResponse {
	caseDTOs := make([]dtos.CaseResponse, len(cases))
	for i, caseModel := range cases {
		caseDTOs[i] = *s.CaseToDTO(&caseModel)
	}
	return caseDTOs
}

func (s *CaseConversionServiceImpl) UpdateCaseFieldsToMap(updateRequest dtos.UpdateCaseRequest) map[string]interface{} {
	updateFields := make(map[string]interface{})

	if updateRequest.Name.Present {
		updateFields["name"] = updateRequest.Name.Value
	}
	if updateRequest.Messages.Present {
		messages, _ := s.DTOToMessages(updateRequest.Messages.Value)
		updateFields["messages"] = messages
	}
	if updateRequest.Collaborators.Present {
		collaborators, _ := s.DTOToCollaborators(updateRequest.Collaborators.Value)
		updateFields["collaborators"] = collaborators
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

	return updateFields
}

func (s *CaseConversionServiceImpl) CollaboratorsToDTO(collaborators []models.Collaborators) []dtos.CollaboratorResponse {
	collaboratorDTOs := make([]dtos.CollaboratorResponse, len(collaborators))
	for i, collaborator := range collaborators {
		collaboratorDTOs[i] = dtos.CollaboratorResponse{
			ID:   helpers.NewNullable(collaborator.ID),
			Edit: helpers.NewNullable(collaborator.Edit),
		}
	}
	return collaboratorDTOs
}

func (s *CaseConversionServiceImpl) DTOToCollaborators(collaboratorsDTO []dtos.CollaboratorResponse) ([]models.Collaborators, error) {
	collaborators := make([]models.Collaborators, len(collaboratorsDTO))
	for i, dto := range collaboratorsDTO {
		if !dto.ID.Present {
			return nil, errors.New("collaborator ID is required")
		}
		collaborators[i] = models.Collaborators{
			ID:   dto.ID.Value,
			Edit: dto.Edit.OrElse(false),
		}
	}
	return collaborators, nil
}

func (s *CaseConversionServiceImpl) MessageToDTO(message models.Message) dtos.MessageResponse {
	return dtos.MessageResponse{
		Content:      helpers.NewNullable(message.Content),
		Sender:       helpers.NewNullable(message.SenderID),
		Recipient:    helpers.NewNullable(message.RecipientID),
		FunctionCall: helpers.NewNullable(message.FunctionCall),
	}
}

func (s *CaseConversionServiceImpl) DTOToMessage(messageDTO dtos.MessageResponse) (models.Message, error) {
	if !messageDTO.Content.Present || !messageDTO.Sender.Present || !messageDTO.Recipient.Present {
		return models.Message{}, errors.New("content, sender, and recipient are required")
	}
	return models.Message{
		Content:      messageDTO.Content.Value,
		SenderID:     messageDTO.Sender.Value,
		RecipientID:  messageDTO.Recipient.Value,
		FunctionCall: messageDTO.FunctionCall.OrElse(false),
	}, nil
}

func (s *CaseConversionServiceImpl) MessagesToDTO(messages []models.Message) []dtos.MessageResponse {
	messageDTOs := make([]dtos.MessageResponse, len(messages))
	for i, message := range messages {
		messageDTOs[i] = s.MessageToDTO(message)
	}
	return messageDTOs
}

func (s *CaseConversionServiceImpl) DTOToMessages(messagesDTO []dtos.MessageResponse) ([]models.Message, error) {
	messages := make([]models.Message, len(messagesDTO))
	for i, dto := range messagesDTO {
		message, err := s.DTOToMessage(dto)
		if err != nil {
			return nil, err
		}
		messages[i] = message
	}
	return messages, nil
}
