package handlers

import (
	"encoding/json"
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
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["id"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	user, err := h.userService.GetUserByID(ctx, id)
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
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["id"]))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var userRequest dtos.UpdateUserRequest
	err = json.NewDecoder(r.Body).Decode(&userRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create a map to store the update fields
	updateFields := make(map[string]interface{})

	// Check each field in the UpdateUserRequest and add non-null values to the updateFields map
	if userRequest.Image.Present {
		updateFields["image"] = userRequest.Image.Value
	}
	if userRequest.Email.Present {
		updateFields["email"] = userRequest.Email.Value
	}
	if userRequest.FirstName.Present {
		updateFields["first_name"] = userRequest.FirstName.Value
	}
	if userRequest.LastName.Present {
		updateFields["last_name"] = userRequest.LastName.Value
	}
	if userRequest.Phone.Present {
		updateFields["phone"] = userRequest.Phone.Value
	}
	if userRequest.CaseIDs.Present {
		updateFields["case_ids"] = userRequest.CaseIDs.Value
	}
	if userRequest.TeamID.Present {
		updateFields["team_id"] = userRequest.TeamID.Value
	}
	if userRequest.AgentIDs.Present {
		updateFields["agent_ids"] = userRequest.AgentIDs.Value
	}
	if userRequest.SubscriptionID.Present {
		updateFields["subscription_id"] = userRequest.SubscriptionID.Value
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
	id, err := primitive.ObjectIDFromHex(strings.TrimSpace(mux.Vars(r)["id"]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = h.userService.DeleteUserByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
