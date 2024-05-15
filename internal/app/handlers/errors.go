package handler

import (
	"encoding/json"
	"net/http"

	http_errors "github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/errors"
)

// ErrorResponse creates a generic error response.
func ErrorResponse(w http.ResponseWriter, err http_errors.SlugError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.HTTPStatus())
	response := map[string]interface{}{
		"error": map[string]string{
			"message": err.Error(),
			"slug":    err.Slug(),
			"type":    err.ErrorType().Type,
		},
	}
	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// NotFoundHandler handles requests for resources that are not found.
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	err := http_errors.NewNotFoundError("Resource not found", "resource-not-found")
	ErrorResponse(w, err)
}

// MethodNotAllowedHandler handles requests with unsupported HTTP methods.
func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	err := http_errors.NewIncorrectInputError("Method not allowed", "method-not-allowed")
	ErrorResponse(w, err)
}

// UnauthorizedHandler handles unauthorized requests.
func UnauthorizedHandler(w http.ResponseWriter, r *http.Request) {
	err := http_errors.NewAuthorizationError("Unauthorized", "unauthorized")
	ErrorResponse(w, err)
}

// BadRequestHandler handles requests with invalid or missing parameters.
func BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	err := http_errors.NewIncorrectInputError("Bad request", "bad-request")
	ErrorResponse(w, err)
}

// InternalServerErrorHandler handles internal server errors.
func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	err := http_errors.NewSlugError("Internal server error", "internal-server-error", http_errors.ErrorTypeUnknown)
	ErrorResponse(w, err)
}

// DatabaseErrorHandler handles database-related errors.
func DatabaseErrorHandler(w http.ResponseWriter, r *http.Request) {
	err := http_errors.NewDatabaseError("Database error", "database-error")
	ErrorResponse(w, err)
}

// TimeoutErrorHandler handles timeout errors.
func TimeoutErrorHandler(w http.ResponseWriter, r *http.Request) {
	err := http_errors.NewTimeoutError("Request timeout", "request-timeout")
	ErrorResponse(w, err)
}

// ConflictErrorHandler handles conflict errors (e.g., duplicate resources).
func ConflictErrorHandler(w http.ResponseWriter, r *http.Request) {
	err := http_errors.NewConflictError("Conflict", "resource-conflict")
	ErrorResponse(w, err)
}

// RateLimitExceededHandler handles rate limit exceeded errors.
func RateLimitExceededHandler(w http.ResponseWriter, r *http.Request) {
	err := http_errors.NewRateLimitError("Rate limit exceeded", "rate-limit-exceeded")
	ErrorResponse(w, err)
}
