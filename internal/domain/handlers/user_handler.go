package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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

// GetUserByID - Handles GET requests for a specific user
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := strings.TrimSpace(mux.Vars(r)["id"])
	userID, err := primitive.ObjectIDFromHex(id)
	fmt.Println(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetUserByEmail - Handles POST requests to retrieve a user by email
func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var emailRequest struct {
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&emailRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := h.userService.GetUserByEmail(ctx, emailRequest.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateUser - Handles POST requests to create a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var userRequest dtos.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createdUser, err := h.userService.CreateUser(ctx, &userRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

// UpdateUser - Handles PUT requests to update a user
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := strings.TrimSpace(mux.Vars(r)["id"])

	var userRequest dtos.UpdateUserRequest
	err := json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a map to store the update fields
	updateFields := make(map[string]interface{})

	// Check each field in the UpdateUserRequest and add non-null values to the updateFields map
	if userRequest.Image.Valid {
		updateFields["image"] = userRequest.Image.Val
	}
	if userRequest.Email.Valid {
		updateFields["email"] = userRequest.Email.Val
	}
	if userRequest.FirstName.Valid {
		updateFields["first_name"] = userRequest.FirstName.Val
	}
	if userRequest.LastName.Valid {
		updateFields["last_name"] = userRequest.LastName.Val
	}
	if userRequest.Phone.Valid {
		updateFields["phone"] = userRequest.Phone.Val
	}
	if userRequest.CaseIDs.Valid {
		updateFields["case_ids"] = userRequest.CaseIDs.Val
	}
	if userRequest.TeamID.Valid {
		updateFields["team_id"] = userRequest.TeamID.Val
	}
	if userRequest.AgentIDs.Valid {
		updateFields["agent_ids"] = userRequest.AgentIDs.Val
	}
	if userRequest.SubscriptionID.Valid {
		updateFields["subscription_id"] = userRequest.SubscriptionID.Val
	}

	// Call the UpdateUser service method with the update fields
	updatedUser, err := h.userService.UpdateUser(ctx, id, updateFields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

// DeleteUser - Handles DELETE requests to delete a user
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := strings.TrimSpace(mux.Vars(r)["id"])
	err := h.userService.DeleteUserByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
