package daos

import (
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CaseDAO interface {
	FindAll() ([]dtos.CaseResponse, error)
	FindByID(id primitive.ObjectID) (dtos.CaseResponse, error)
	FindByCreatorID(creatorID primitive.ObjectID) ([]dtos.CaseResponse, error)
	Create(caseRequest dtos.CreateCaseRequest) error
	Update(id primitive.ObjectID, updates dtos.UpdateCaseRequest) error
	Delete(id primitive.ObjectID) error
	AddCollaborator(id primitive.ObjectID, collaboratorID primitive.ObjectID) error
	RemoveCollaborator(id primitive.ObjectID, collaboratorID primitive.ObjectID) error
}
