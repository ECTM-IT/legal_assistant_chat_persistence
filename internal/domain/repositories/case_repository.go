package repositories

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CaseRepository struct {
	caseDAO *daos.CaseDAO
}

func NewCaseRepository(caseDAO *daos.CaseDAO) *CaseRepository {
	return &CaseRepository{
		caseDAO: caseDAO,
	}
}

func (r *CaseRepository) GetAllCases() ([]dtos.CaseResponse, error) {
	return r.caseDAO.FindAll()
}

func (r *CaseRepository) GetCaseByID(id primitive.ObjectID) (dtos.CaseResponse, error) {
	return r.caseDAO.FindByID(id)
}

func (r *CaseRepository) GetCasesByCreatorID(creatorID string) ([]dtos.CaseResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(creatorID)
	if err != nil {
		return nil, err
	}
	return r.caseDAO.FindByCreatorID(objectID)
}

func (r *CaseRepository) CreateCase(caseRequest dtos.CreateCaseRequest) error {
	return r.caseDAO.Create(caseRequest)
}

func (r *CaseRepository) UpdateCase(id primitive.ObjectID, updates dtos.UpdateCaseRequest) error {
	return r.caseDAO.Update(id, updates)
}

func (r *CaseRepository) DeleteCase(id primitive.ObjectID) error {
	return r.caseDAO.Delete(id)
}

func (r *CaseRepository) AddCollaboratorToCase(id primitive.ObjectID, collaboratorID primitive.ObjectID) error {
	return r.caseDAO.AddCollaborator(id, collaboratorID)
}

func (r *CaseRepository) RemoveCollaboratorFromCase(id primitive.ObjectID, collaboratorID primitive.ObjectID) error {
	return r.caseDAO.RemoveCollaborator(id, collaboratorID)
}
