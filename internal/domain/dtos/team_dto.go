package dtos

import (
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateTeamRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type TeamResponse struct {
	ID          primitive.ObjectID   `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Members     []TeamMemberResponse `json:"members"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

type UpdateTeamRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type TeamMemberResponse struct {
	ID         primitive.ObjectID `json:"id"`
	UserID     primitive.ObjectID `json:"user_id"`
	Role       models.Role        `json:"role"`
	FirstName  string             `json:"first_name"`
	LastName   string             `json:"last_name"`
	Email      string             `json:"email"`
	DateAdded  time.Time          `json:"date_added"`
	LastActive time.Time          `json:"last_active"`
}

type AddTeamMemberRequest struct {
	Email     string      `json:"email" validate:"required,email"`
	Role      models.Role `json:"role" validate:"required"`
	FirstName string      `json:"first_name" validate:"required"`
	LastName  string      `json:"last_name" validate:"required"`
}

type UpdateTeamMemberRequest struct {
	Role      *models.Role `json:"role,omitempty"`
	FirstName *string      `json:"first_name,omitempty"`
	LastName  *string      `json:"last_name,omitempty"`
	Email     *string      `json:"email,omitempty" validate:"omitempty,email"`
}

type TeamInvitationRequest struct {
	Email string      `json:"email" validate:"required,email"`
	Role  models.Role `json:"role" validate:"required"`
}

type TeamInvitationResponse struct {
	ID        primitive.ObjectID `json:"id"`
	Email     string             `json:"email"`
	Role      models.Role        `json:"role"`
	CreatedAt time.Time          `json:"created_at"`
	ExpiresAt time.Time          `json:"expires_at"`
	IsUsed    bool               `json:"is_used"`
}

type AcceptInvitationRequest struct {
	Token     string `json:"token" validate:"required"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}
