package handler

import (
	"encoding/json"
	"net/http"

	http_errors "github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/errors"
)

// errorResponse is a helper function to create a generic error response
func errorResponse(w http.ResponseWriter, err http_errors.SlugError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.HTTPStatus())
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"message": err.Error(),
			"slug":    err.Slug(),
			"type":    err.ErrorType().T,
		},
	})
}

// NotFoundHandler handles requests for resources that are not found
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	err := http_errors.NewNotFoundError("Resource not found", "resource-not-found")
	errorResponse(w, err)
}

// MethodNotAllowedHandler handles requests with unsupported HTTP methods
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	err := http_errors.NewIncorrectInputError("Method not allowed", "method-not-allowed")
	errorResponse(w, err)
}

// UnauthorizedHandler handles unauthorized requests
func UnauthorizedHandler(w http.ResponseWriter, r *http.Request) {
	err := http_errors.NewAuthorizationError("Unauthorized", "unauthorized")
	errorResponse(w, err)
}

// BadRequestHandler handles requests with invalid or missing parameters
func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	err := http_errors.NewIncorrectInputError("Bad request", "bad-request")
	errorResponse(w, err)
}

// InternalServerErrorHandler handles internal server errors
func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	err := http_errors.NewSlugError("Internal server error", "internal-server-error", http_errors.ErrorTypeUnknown)
	errorResponse(w, err)
}

// DatabaseErrorHandler handles database-related errors
func DatabaseErrorHandler(w http.ResponseWriter, r *http.Request) {
	err := http_errors.NewDatabaseError("Database error", "database-error")
	errorResponse(w, err)
}

// TimeoutErrorHandler handles timeout errors
func TimeoutErrorHandler(w http.ResponseWriter, r *http.Request) {
	err := http_errors.NewTimeoutError("Request timeout", "request-timeout")
	errorResponse(w, err)
}

// ConflictErrorHandler handles conflict errors (e.g., duplicate resources)
func ConflictErrorHandler(w http.ResponseWriter, r *http.Request) {
	err := http_errors.NewConflictError("Conflict", "resource-conflict")
	errorResponse(w, err)
}

// RateLimitExceededHandler handles rate limit exceeded errors
func RateLimitExceededHandler(w http.ResponseWriter, r *http.Request) {
	err := http_errors.NewRateLimitError("Rate limit exceeded", "rate-limit-exceeded")
	errorResponse(w, err)
}
