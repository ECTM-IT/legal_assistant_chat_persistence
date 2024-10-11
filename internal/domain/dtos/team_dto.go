package dtos

import (
	"time"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/app/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateTeamRequest struct {
	AdminID helpers.Nullable[primitive.ObjectID]   `json:"admin_id" bson:"admin_id"`
	Members helpers.Nullable[[]TeamMemberResponse] `json:"members" bson:"members"`
}

type TeamResponse struct {
	ID      helpers.Nullable[primitive.ObjectID]   `json:"id" bson:"_id,omitempty"`
	AdminID helpers.Nullable[primitive.ObjectID]   `json:"admin_id" bson:"admin_id"`
	Members helpers.Nullable[[]TeamMemberResponse] `json:"members" bson:"members"`
}

type UpdateTeamRequest struct {
	AdminID helpers.Nullable[primitive.ObjectID]   `json:"admin_id" bson:"admin_id,omitempty"`
	Members helpers.Nullable[[]TeamMemberResponse] `json:"members" bson:"members,omitempty"`
}

type TeamMemberResponse struct {
	ID         helpers.Nullable[primitive.ObjectID] `json:"id" bson:"_id,omitempty"`
	UserID     helpers.Nullable[primitive.ObjectID] `json:"user_id" bson:"user_id"`
	DateAdded  helpers.Nullable[time.Time]          `json:"date_added" bson:"date_added"`
	LastActive helpers.Nullable[time.Time]          `json:"last_active" bson:"last_active"`
}

type AddMemberRequest struct {
	Email helpers.Nullable[string] `json:"email" bson:"email"`
}

type ChangeAdminRequest struct {
	Email helpers.Nullable[string] `json:"email" bson:"email"`
}
