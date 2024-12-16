package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	RoleAdmin Role = "Admin User"
	RoleUser  Role = "Account User"
)

type Team struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description,omitempty"`
	Members     []TeamMember       `bson:"members"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
	IsDeleted   bool               `bson:"is_deleted"`
	DeletedAt   *time.Time         `bson:"deleted_at,omitempty"`
}

type TeamMember struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UserID     primitive.ObjectID `bson:"user_id"`
	Role       Role               `bson:"role"`
	FirstName  string             `bson:"first_name"`
	LastName   string             `bson:"last_name"`
	Email      string             `bson:"email"`
	DateAdded  time.Time          `bson:"date_added"`
	LastActive time.Time          `bson:"last_active"`
	IsDeleted  bool               `bson:"is_deleted"`
	DeletedAt  *time.Time         `bson:"deleted_at,omitempty"`
}

type TeamInvitation struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	TeamID    primitive.ObjectID `bson:"team_id"`
	Email     string             `bson:"email"`
	Role      Role               `bson:"role"`
	Token     string             `bson:"token"`
	CreatedAt time.Time          `bson:"created_at"`
	ExpiresAt time.Time          `bson:"expires_at"`
	IsUsed    bool               `bson:"is_used"`
}
