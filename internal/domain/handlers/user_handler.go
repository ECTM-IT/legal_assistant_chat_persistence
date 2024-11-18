package handlers

import (
	"context"
	"net/http"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*dtos.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*dtos.UserResponse, error)
	CreateUser(ctx context.Context, req *dtos.CreateUserRequest) (*dtos.UserResponse, error)
	UpdateUser(ctx context.Context, id primitive.ObjectID, updateFields map[string]interface{}) (*dtos.UserResponse, error)
	DeleteUserByID(ctx context.Context, id primitive.ObjectID) error
}

type UserHandler struct {
	BaseHandler
	service *services.UserServiceImpl
}

func NewUserHandler(service *services.UserServiceImpl) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := h.service.GetUserByID(r.Context(), id)
	if err != nil {
		h.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	var emailRequest struct {
		Email string `json:"email"`
	}
	if err := h.DecodeJSONBody(r, &emailRequest); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.service.GetUserByEmail(r.Context(), emailRequest.Email)
	if err != nil {
		h.RespondWithError(w, http.StatusNotFound, "User not found")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dtos.CreateUserRequest
	if err := h.DecodeJSONBody(r, &req); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	createdUser, err := h.service.CreateUser(r.Context(), &req)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	h.RespondWithJSON(w, http.StatusCreated, createdUser)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req dtos.UpdateUserRequest
	if err := h.DecodeJSONBody(r, &req); err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	updatedUser, err := h.service.UpdateUser(r.Context(), id, &req)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}
	h.RespondWithJSON(w, http.StatusOK, updatedUser)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := h.ParseObjectID(r, "id", false)
	if err != nil {
		h.RespondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = h.service.DeleteUserByID(r.Context(), id)
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}
	h.RespondWithJSON(w, http.StatusNoContent, nil)
}
