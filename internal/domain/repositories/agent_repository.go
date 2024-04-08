package repositories

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/daos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
)

type AgentRepository struct {
	agentDAO *daos.AgentDAO
	userDAO  *daos.UserDAO
}

func NewAgentRepository(agentDAO *daos.AgentDAO, userDAO *daos.UserDAO) *AgentRepository {
	return &AgentRepository{
		agentDAO: agentDAO,
		userDAO:  userDAO,
	}
}

func (r *AgentRepository) GetAllAgents(ctx context.Context) ([]dtos.AgentResponse, error) {
	agents, err := r.agentDAO.GetAllAgents(ctx)
	if err != nil {
		return nil, err
	}
	return r.toAgentResponses(agents), nil
}

func (r *AgentRepository) GetAgentByID(ctx context.Context, id string) (*dtos.AgentResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	agent, err := r.agentDAO.GetAgentByID(ctx, objectID)
	if err != nil {
		return nil, err
	}
	return r.toAgentResponse(agent), nil
}

func (r *AgentRepository) GetAgentsByUserID(ctx context.Context, userID primitive.ObjectID) ([]dtos.AgentResponse, error) {
	user, err := r.userDAO.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	agents, err := r.agentDAO.GetAgentsByIDs(ctx, user.AgentIDs)
	if err != nil {
		return nil, err
	}
	return r.toAgentResponses(agents), nil
}

func (r *AgentRepository) PurchaseAgent(ctx context.Context, userID primitive.ObjectID, agentID primitive.ObjectID) (*dtos.UserResponse, error) {
	user, err := r.userDAO.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if contains(user.AgentIDs, agentID) {
		return nil, errors.New("agent already added to the user")
	}
	_, err = r.agentDAO.GetAgentByID(ctx, agentID)
	if err != nil {
		return nil, err
	}
	user.AgentIDs = append(user.AgentIDs, agentID)
	err = r.userDAO.UpdateUser(userID, user)
	if err != nil {
		return nil, err
	}
	return r.toUserResponse(user), nil
}

func (r *AgentRepository) toAgentResponse(agent *models.Agent) *dtos.AgentResponse {
	var skillResponses []dtos.SkillResponse
	for _, skill := range agent.Skills {
		skillResponses = append(skillResponses, dtos.SkillResponse{
			Name:         helpers.NewNullable(skill.Name),
			Descriptions: helpers.NewNullable(skill.Descriptions),
		})
	}
	return &dtos.AgentResponse{
		ID:           helpers.NewNullable(agent.ID),
		ProfileImage: helpers.NewNullable(agent.ProfileImage),
		Name:         helpers.NewNullable(agent.Name),
		Description:  helpers.NewNullable(agent.Description),
		Skills:       helpers.NewNullable(skillResponses),
		Price:        helpers.NewNullable(agent.Price),
		Code:         helpers.NewNullable(agent.Code),
	}
}

func (r *AgentRepository) toAgentResponses(agents []models.Agent) []dtos.AgentResponse {
	var agentResponses []dtos.AgentResponse
	for _, agent := range agents {
		agentResponses = append(agentResponses, *r.toAgentResponse(&agent))
	}
	return agentResponses
}

func (r *AgentRepository) toUserResponse(user *models.User) *dtos.UserResponse {
	return &dtos.UserResponse{
		ID:             helpers.NewNullable(user.ID),
		Image:          helpers.NewNullable(user.Image),
		Email:          helpers.NewNullable(user.Email),
		FirstName:      helpers.NewNullable(user.FirstName),
		LastName:       helpers.NewNullable(user.LastName),
		Phone:          helpers.NewNullable(user.Phone),
		CaseIDs:        helpers.NewNullable(user.CaseIDs),
		TeamID:         helpers.NewNullable(user.TeamID),
		AgentIDs:       helpers.NewNullable(user.AgentIDs),
		SubscriptionID: helpers.NewNullable(user.SubscriptionID),
	}
}

func contains(slice []primitive.ObjectID, element primitive.ObjectID) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}
