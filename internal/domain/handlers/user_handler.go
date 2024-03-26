package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/dtos"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/domain/services"
)

// UserHandler - Handles user-related HTTP requests
type UserHandler struct {
	userService *services.UserServiceImpl
}

// NewUserHandler - Creates a new UserHandler
func NewUserHandler(userService *services.UserServiceImpl) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUser - Handles GET requests for a specific user
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := primitive.ObjectIDFromHex(vars["userID"])
	if err != nil {
		// Handle error appropriately (e.g., write specific error responses)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		// Handle error appropriately (e.g., write specific error responses)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateUser - Handles POST requests to create a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userRequest dtos.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdUser, err := h.userService.CreateUser(&userRequest)
	if err != nil {
		// Handle error appropriately
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}
