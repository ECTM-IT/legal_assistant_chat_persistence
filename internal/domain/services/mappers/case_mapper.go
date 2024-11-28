// Start of Selection
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
	UpdateCaseFieldsToMap(updateRequest dtos.UpdateCaseRequest) (map[string]interface{}, error)
	CollaboratorsToDTO(collaborators []models.Collaborators) []dtos.CollaboratorResponse
	DTOToCollaborators(collaboratorsDTO []dtos.CollaboratorResponse) ([]models.Collaborators, error)
	MessageToDTO(message models.Message) dtos.MessageResponse
	DTOToMessage(messageDTO dtos.MessageResponse) (models.Message, error)
	MessagesToDTO(messages []models.Message) []dtos.MessageResponse
	DTOToMessages(messagesDTO []dtos.MessageResponse) ([]models.Message, error)
	DocumentsToDocumentsResponse(docs []models.Document) []dtos.DocumentResponse
	DTOToDocuments(documentsDTO []dtos.DocumentResponse) ([]models.Document, error)
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
		err := errors.New("creator ID is required")
		s.logger.Error("Failed to convert DTO to Case: creator ID is required", err)
		return nil, err
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

	documents, err := s.DTOToDocuments(caseRequest.Documents.OrElse(nil))
	if err != nil {
		s.logger.Error("Failed to convert documents", err)
		return nil, fmt.Errorf("error converting documents: %w", err)
	}

	now := time.Now()
	caseModel := &models.Case{
		ID:            primitive.NewObjectID(),
		Name:          caseRequest.Name.OrElse("New Case"),
		CreatorID:     caseRequest.CreatorID.Value,
		Messages:      messages,
		Collaborators: collaborators,
		Documents:     documents,
		Skills:        caseRequest.Skills.OrElse([]string{}),
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
		Documents:     helpers.NewNullable(s.DocumentsToDTO(caseModel.Documents)),
		Skills:        helpers.NewNullable(caseModel.Skills),
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

	if len(cases) == 0 {
		s.logger.Warn("No Cases provided for conversion")
		return []dtos.CaseResponse{}
	}

	caseDTOs := make([]dtos.CaseResponse, 0, len(cases))
	for _, caseModel := range cases {
		dto := s.CaseToDTO(&caseModel)
		if dto != nil {
			caseDTOs = append(caseDTOs, *dto)
		}
	}

	s.logger.Info("Successfully converted multiple Cases to DTOs")
	return caseDTOs
}

func (s *CaseConversionServiceImpl) UpdateCaseFieldsToMap(updateRequest dtos.UpdateCaseRequest) (map[string]interface{}, error) {
	s.logger.Info("Converting UpdateCaseRequest to map")
	updateFields := make(map[string]interface{})

	if updateRequest.Name.Present {
		updateFields["name"] = updateRequest.Name.Value
	}
	if updateRequest.Messages.Present {
		messages, err := s.DTOToMessages(updateRequest.Messages.Value)
		if err != nil {
			s.logger.Error("Failed to convert messages for update", err)
			return nil, fmt.Errorf("error converting messages: %w", err)
		}
		updateFields["messages"] = messages
	}
	if updateRequest.Collaborators.Present {
		collaborators, err := s.DTOToCollaborators(updateRequest.Collaborators.Value)
		if err != nil {
			s.logger.Error("Failed to convert collaborators for update", err)
			return nil, fmt.Errorf("error converting collaborators: %w", err)
		}
		updateFields["collaborators"] = collaborators
	}
	if updateRequest.Documents.Present {
		documents, err := s.DTOToDocuments(updateRequest.Documents.Value)
		if err != nil {
			s.logger.Error("Failed to convert documents for update", err)
			return nil, fmt.Errorf("error converting documents: %w", err)
		}
		updateFields["documents"] = documents
	}
	if updateRequest.Skills.Present {
		updateFields["skills"] = updateRequest.Skills.Value
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
	return updateFields, nil
}

func (s *CaseConversionServiceImpl) CollaboratorsToDTO(collaborators []models.Collaborators) []dtos.CollaboratorResponse {
	s.logger.Info("Converting Collaborators to DTOs")

	if len(collaborators) == 0 {
		s.logger.Warn("No Collaborators provided for conversion")
		return []dtos.CollaboratorResponse{}
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

	if len(collaboratorsDTO) == 0 {
		s.logger.Warn("No CollaboratorResponses provided for conversion")
		return []models.Collaborators{}, nil
	}

	collaborators := make([]models.Collaborators, 0, len(collaboratorsDTO))
	for _, dto := range collaboratorsDTO {
		if !dto.ID.Present {
			err := errors.New("collaborator ID is required")
			s.logger.Error("Failed to convert DTO to Collaborator: collaborator ID is required", err)
			return nil, err
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

	feedbacks := make([]dtos.Feedback, len(message.Feedbacks))
	for _, f := range message.Feedbacks {
		feedbacks = append(feedbacks, dtos.Feedback{
			ID:           helpers.NewNullable(f.ID),
			CaseID:       helpers.NewNullable(f.CaseID),
			MessageID:    helpers.NewNullable(f.MessageID),
			CreatorID:    helpers.NewNullable(f.CreatorID),
			Score:        helpers.NewNullable(f.Score),
			Reasons:      helpers.NewNullable(f.Reasons),
			Comment:      helpers.NewNullable(f.Comment),
			CreationDate: helpers.NewNullable(f.CreationDate),
		})
	}

	dto := dtos.MessageResponse{
		ID:           helpers.NewNullable(message.ID),
		Content:      helpers.NewNullable(message.Content),
		Sender:       helpers.NewNullable(message.Sender),
		Recipient:    helpers.NewNullable(message.Recipient),
		FunctionCall: helpers.NewNullable(message.FunctionCall),
		DocumentPath: helpers.NewNullable(message.DocumentPath),
		Feedbacks:    helpers.NewNullable(feedbacks),

		Skills: helpers.NewNullable(message.Skills),
		Agent:  helpers.NewNullable(message.Agent),
	}

	s.logger.Info("Successfully converted Message to DTO")
	return dto
}

func (s *CaseConversionServiceImpl) DTOToMessage(messageDTO dtos.MessageResponse) (models.Message, error) {
	s.logger.Info("Converting DTO to Message")

	s.logger.Info(messageDTO.Content.String() + messageDTO.Recipient.String() + messageDTO.Sender.String())

	if !messageDTO.Content.Present || !messageDTO.Sender.Present || !messageDTO.Recipient.Present {
		err := errors.New("content, sender, and recipient are required")
		s.logger.Error("Failed to convert DTO to Message: content, sender, and recipient are required", err)
		return models.Message{}, err
	}

	feedbacks := make([]models.Feedback, len(messageDTO.Feedbacks.Value))
	for _, f := range messageDTO.Feedbacks.Value {
		feedbacks = append(feedbacks, models.Feedback{
			ID:           f.ID.Value,
			CaseID:       f.CaseID.Value,
			MessageID:    f.MessageID.Value,
			CreatorID:    f.CreatorID.Value,
			Score:        f.Score.Value,
			Reasons:      f.Reasons.Value,
			Comment:      f.Comment.Value,
			CreationDate: f.CreationDate.Value,
		})
	}

	message := models.Message{
		ID:           messageDTO.ID.Value,
		Content:      messageDTO.Content.Value,
		Sender:       messageDTO.Sender.Value,
		Recipient:    messageDTO.Recipient.Value,
		FunctionCall: messageDTO.FunctionCall.OrElse(false),
		DocumentPath: messageDTO.DocumentPath.OrElse(""),
		Feedbacks:    feedbacks,

		Skills: messageDTO.Skills.OrElse([]string{}),
		Agent:  messageDTO.Agent.OrElse(""),
	}

	s.logger.Info("Successfully converted DTO to Message")
	return message, nil
}

func (s *CaseConversionServiceImpl) MessagesToDTO(messages []models.Message) []dtos.MessageResponse {
	s.logger.Info("Converting multiple Messages to DTOs")

	if len(messages) == 0 {
		s.logger.Warn("No Messages provided for conversion")
		return []dtos.MessageResponse{}
	}

	messageDTOs := make([]dtos.MessageResponse, 0, len(messages))
	for _, message := range messages {
		messageDTOs = append(messageDTOs, s.MessageToDTO(message))
	}

	s.logger.Info("Successfully converted multiple Messages to DTOs")
	return messageDTOs
}

func (s *CaseConversionServiceImpl) DTOToMessages(messagesDTO []dtos.MessageResponse) ([]models.Message, error) {
	s.logger.Info("Converting DTOs to Messages")

	if len(messagesDTO) == 0 {
		s.logger.Warn("No MessageResponses provided for conversion")
		return []models.Message{}, nil
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

// DocumentsToDTO converts a list of Document models to DocumentResponse DTOs
func (s *CaseConversionServiceImpl) DocumentsToDTO(docs []models.Document) []dtos.DocumentResponse {
	s.logger.Info("Converting Documents to DTOs")

	var docResponses []dtos.DocumentResponse
	for _, doc := range docs {
		docResponses = append(docResponses, dtos.DocumentResponse{
			ID:          helpers.NewNullable(doc.ID),
			FileName:    helpers.NewNullable(doc.FileName),
			FileType:    helpers.NewNullable(doc.FileType),
			FileContent: helpers.NewNullable(doc.FileContent),
			UploadDate:  helpers.NewNullable(doc.UploadDate),
		})
	}

	s.logger.Info("Successfully converted Documents to DTOs")
	return docResponses
}

func (s *CaseConversionServiceImpl) DTOToDocuments(documentsDTO []dtos.DocumentResponse) ([]models.Document, error) {
	s.logger.Info("Converting DTOs to Documents")

	if len(documentsDTO) == 0 {
		s.logger.Warn("No DocumentResponses provided for conversion")
		return []models.Document{}, nil
	}

	documents := make([]models.Document, 0, len(documentsDTO))
	for _, dto := range documentsDTO {
		if !dto.ID.Present {
			err := errors.New("document ID is required")
			s.logger.Error("Failed to convert DTO to Document: document ID is required", err)
			return nil, err
		}
		documents = append(documents, models.Document{
			ID:          dto.ID.Value,
			FileName:    dto.FileName.OrElse(""),
			FileType:    dto.FileType.OrElse(""),
			FileContent: dto.FileContent.OrElse([]byte{}),
			UploadDate:  dto.UploadDate.OrElse(time.Now()),
		})
	}

	s.logger.Info("Successfully converted DTOs to Documents")
	return documents, nil
}
