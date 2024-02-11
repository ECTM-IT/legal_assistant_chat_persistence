package handler

import (
	"encoding/json"
	"net/http"

	http_errors "github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/errors"
)

// Function to create a generic error response
func errorResponse(w http.ResponseWriter, err http_errors.SlugError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.HTTPStatus())
	json.NewEncoder(w).Encode(map[string]string{
		"error": err.Error(),
		"slug":  err.Slug(),
	})
}

// Specific error handling functions can be created as follows:
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	err := http_errors.NewNotFoundError("Resource not found", "not-found")
	errorResponse(w, err)
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	err := http_errors.NewSlugError("Method not allowed", "method-not-allowed", http_errors.ErrorTypeIncorrectInput) // Assuming you have defined this constructor
	errorResponse(w, err)
}
