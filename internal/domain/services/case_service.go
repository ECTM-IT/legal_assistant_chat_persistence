package services

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CaseService struct {
	caseDAO daos.CaseDAO
}

func NewCaseService(caseDAO daos.CaseDAO) *CaseService {
	return &CaseService{
		caseDAO: caseDAO,
	}
}

func (s *CaseService) GetAllCases() ([]dtos.CaseResponse, error) {
	return s.caseDAO.FindAll()
}

func (s *CaseService) GetCaseByID(id string) (dtos.CaseResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return dtos.CaseResponse{}, err
	}
	return s.caseDAO.FindByID(objectID)
}

func (s *CaseService) GetCasesByCreatorID(creatorID string) ([]dtos.CaseResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(creatorID)
	if err != nil {
		return nil, err
	}
	return s.caseDAO.FindByCreatorID(objectID)
}

func (s *CaseService) CreateCase(caseRequest dtos.CreateCaseRequest) error {
	return s.caseDAO.Create(caseRequest)
}

func (s *CaseService) UpdateCase(id string, updates dtos.UpdateCaseRequest) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.caseDAO.Update(objectID, updates)
}

func (s *CaseService) DeleteCase(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.caseDAO.Delete(objectID)
}

func (s *CaseService) AddCollaboratorToCase(id, collaboratorID string) error {
	caseID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	collaboratorObjectID, err := primitive.ObjectIDFromHex(collaboratorID)
	if err != nil {
		return err
	}
	return s.caseDAO.AddCollaborator(caseID, collaboratorObjectID)
}

func (s *CaseService) RemoveCollaboratorFromCase(id, collaboratorID string) error {
	caseID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	collaboratorObjectID, err := primitive.ObjectIDFromHex(collaboratorID)
	if err != nil {
		return err
	}
	return s.caseDAO.RemoveCollaborator(caseID, collaboratorObjectID)
}
