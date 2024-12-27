package dtos

import (
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateTeamRequest struct {
	Name        helpers.Nullable[string]              `json:"name" validate:"required"`
	Description helpers.Nullable[string]              `json:"description"`
	Members     helpers.Nullable[[]TeamMemberRequest] `json:"members"`
}

type TeamResponse struct {
	ID          helpers.Nullable[primitive.ObjectID]   `json:"id"`
	Name        helpers.Nullable[string]               `json:"name"`
	Description helpers.Nullable[string]               `json:"description"`
	Members     helpers.Nullable[[]TeamMemberResponse] `json:"members"`
	CreatedAt   helpers.Nullable[time.Time]            `json:"created_at"`
	UpdatedAt   helpers.Nullable[time.Time]            `json:"updated_at"`
}

type UpdateTeamRequest struct {
	Name        helpers.Nullable[string] `json:"name,omitempty"`
	Description helpers.Nullable[string] `json:"description,omitempty"`
}

type TeamMemberRequest struct {
	Email     helpers.Nullable[string]             `json:"email,omitempty"`
	UserID    helpers.Nullable[primitive.ObjectID] `json:"user_id,omitempty"`
	Role      helpers.Nullable[models.Role]        `json:"role,omitempty"`
	FirstName helpers.Nullable[string]             `json:"first_name,omitempty"`
	LastName  helpers.Nullable[string]             `json:"last_name,omitempty"`
}

type TeamMemberResponse struct {
	ID         helpers.Nullable[primitive.ObjectID] `json:"id"`
	UserID     helpers.Nullable[primitive.ObjectID] `json:"user_id"`
	Role       helpers.Nullable[models.Role]        `json:"role"`
	FirstName  helpers.Nullable[string]             `json:"first_name"`
	LastName   helpers.Nullable[string]             `json:"last_name"`
	Email      helpers.Nullable[string]             `json:"email"`
	DateAdded  helpers.Nullable[time.Time]          `json:"date_added"`
	LastActive helpers.Nullable[time.Time]          `json:"last_active"`
}

type AddTeamMemberRequest struct {
	Email     helpers.Nullable[string]      `json:"email" validate:"required,email"`
	Role      helpers.Nullable[models.Role] `json:"role" validate:"required"`
	FirstName helpers.Nullable[string]      `json:"first_name" validate:"required"`
	LastName  helpers.Nullable[string]      `json:"last_name" validate:"required"`
}

type UpdateTeamMemberRequest struct {
	Role      helpers.Nullable[models.Role] `json:"role,omitempty"`
	FirstName helpers.Nullable[string]      `json:"first_name,omitempty"`
	LastName  helpers.Nullable[string]      `json:"last_name,omitempty"`
	Email     helpers.Nullable[string]      `json:"email,omitempty" validate:"omitempty,email"`
}

type TeamInvitationRequest struct {
	Email helpers.Nullable[string]      `json:"email" validate:"required,email"`
	Role  helpers.Nullable[models.Role] `json:"role" validate:"required"`
}

type TeamInvitationResponse struct {
	ID        helpers.Nullable[primitive.ObjectID] `json:"id"`
	Email     helpers.Nullable[string]             `json:"email"`
	Role      helpers.Nullable[models.Role]        `json:"role"`
	CreatedAt helpers.Nullable[time.Time]          `json:"created_at"`
	ExpiresAt helpers.Nullable[time.Time]          `json:"expires_at"`
	IsUsed    helpers.Nullable[bool]               `json:"is_used"`
}

type AcceptInvitationRequest struct {
	Token     helpers.Nullable[string] `json:"token" validate:"required"`
	FirstName helpers.Nullable[string] `json:"first_name" validate:"required"`
	LastName  helpers.Nullable[string] `json:"last_name" validate:"required"`
}
