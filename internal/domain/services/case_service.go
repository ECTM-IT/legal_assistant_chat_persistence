package services

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CaseService struct {
	caseRepo *repositories.CaseRepository
}

func NewCaseService(caseRepo *repositories.CaseRepository) *CaseService {
	return &CaseService{
		caseRepo: caseRepo,
	}
}

func (s *CaseService) GetAllCases() ([]dtos.CaseResponse, error) {
	return s.caseRepo.GetAllCases()
}

func (s *CaseService) GetCaseByID(id string) (dtos.CaseResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return dtos.CaseResponse{}, err
	}
	return s.caseRepo.GetCaseByID(objectID)
}

func (s *CaseService) GetCasesByCreatorID(creatorID string) ([]dtos.CaseResponse, error) {
	return s.caseRepo.GetCasesByCreatorID(creatorID)
}

func (s *CaseService) CreateCase(caseRequest dtos.CreateCaseRequest) error {
	return s.caseRepo.CreateCase(caseRequest)
}

func (s *CaseService) UpdateCase(id primitive.ObjectID, updates dtos.UpdateCaseRequest) error {
	return s.caseRepo.UpdateCase(id, updates)
}

func (s *CaseService) DeleteCase(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return s.caseRepo.DeleteCase(objectID)
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
	return s.caseRepo.AddCollaboratorToCase(caseID, collaboratorObjectID)
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
	return s.caseRepo.RemoveCollaboratorFromCase(caseID, collaboratorObjectID)
}
