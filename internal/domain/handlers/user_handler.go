package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	dto "github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	services "github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
)

// User - A simplified user representation
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserHandler - Handles user-related HTTP requests
type UserHandler struct {
	userService services.UserService
}

// NewUserHandler - Creates a new UserHandler
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUser - Handles GET requests for a specific user
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	user, err := h.userService.GetUserByUserID(ctx, userID)
	if err != nil {
		// Handle error appropriately (e.g., write specific error responses)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateUser - Handles POST requests to create a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request, user *dto.CreateUserRequest) {
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.userService.CreateUser(r.Context(), user)
	if err != nil {
		// Handle error appropriately
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}
