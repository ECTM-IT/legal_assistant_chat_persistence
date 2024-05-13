package dtos

import (
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateTeamRequest struct {
	ID      helpers.Nullable[primitive.ObjectID]   `json:"id"`
	AdminID helpers.Nullable[primitive.ObjectID]   `json:"admin_id"`
	Members helpers.Nullable[[]TeamMemberResponse] `json:"members"`
}
type TeamResponse struct {
	ID      helpers.Nullable[primitive.ObjectID]   `json:"id"`
	AdminID helpers.Nullable[primitive.ObjectID]   `json:"admin_id"`
	Members helpers.Nullable[[]TeamMemberResponse] `json:"members"`
}

type UpdateTeamRequest struct {
	ID      helpers.Nullable[primitive.ObjectID]   `json:"id"`
	AdminID helpers.Nullable[primitive.ObjectID]   `json:"admin_id"`
	Members helpers.Nullable[[]TeamMemberResponse] `json:"members"`
}

type TeamMemberResponse struct {
	ID         helpers.Nullable[primitive.ObjectID] `json:"id"`
	UserID     helpers.Nullable[primitive.ObjectID] `json:"user_id"`
	DateAdded  helpers.Nullable[time.Time]          `json:"date_added"`
	LastActive helpers.Nullable[time.Time]          `json:"last_active"`
}

type AddMemberRequest struct {
	Email helpers.Nullable[string] `json:"email"`
}

type ChangeAdminRequest struct {
	Email helpers.Nullable[string] `json:"email"`
}
