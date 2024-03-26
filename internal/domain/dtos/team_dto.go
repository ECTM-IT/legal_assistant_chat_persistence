package dtos

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeamResponse struct {
	ID      primitive.ObjectID   `json:"id"`
	AdminID primitive.ObjectID   `json:"admin_id"`
	Members []TeamMemberResponse `json:"members"`
}

type TeamMemberResponse struct {
	ID         primitive.ObjectID `json:"id"`
	UserID     primitive.ObjectID `json:"user_id"`
	DateAdded  time.Time          `json:"date_added"`
	LastActive time.Time          `json:"last_active"`
}

type AddMemberRequest struct {
	Email string `json:"email"`
}

type ChangeAdminRequest struct {
	Email string `json:"email"`
}
