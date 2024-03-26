package errors

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
	"go.uber.org/zap"
)

// ErrorResponse struct for sending error responses
type ErrorResponse struct {
	Slug       string `json:"slug"`
	HttpStatus int    `json:"-"`
}

// DefaultErrorResponder contains methods for responding to different types of errors
type DefaultErrorResponder struct {
	Logger *logs.ZapLogger
}

// RespondWithError provides response handling logic
func (responder *DefaultErrorResponder) RespondWithError(ctx context.Context, err error, slug string, w http.ResponseWriter, r *http.Request) {
	// Example error type checking, replace errors.SlugError with your actual error types
	if slugError, ok := err.(SlugError); ok { // SlugError type needs to be defined elsewhere
		responder.handleSlugError(ctx, slugError, w, r)
		return
	}
	// Default to internal server error
	responder.internalServerError(ctx, err, slug, w, r)
}

// Assuming SlugError is an error type you have that includes a Slug() and HTTPStatus() method
func (responder *DefaultErrorResponder) handleSlugError(ctx context.Context, slugError SlugError, w http.ResponseWriter, r *http.Request) {
	status := slugError.HTTPStatus()
	responder.httpRespondWithError(ctx, slugError, slugError.Slug(), w, r, slugError.Error(), status)
}

func (responder *DefaultErrorResponder) internalServerError(ctx context.Context, err error, slug string, w http.ResponseWriter, r *http.Request) {
	responder.httpRespondWithError(ctx, err, slug, w, r, "Internal server error", http.StatusInternalServerError)
}

func (responder *DefaultErrorResponder) httpRespondWithError(_ context.Context, err error, slug string, w http.ResponseWriter, _ *http.Request, logMsg string, status int) {
	logger := responder.Logger // Using the embedded ZapLogger
	if logger != nil {
		logger.Error(logMsg, err, zap.String("error-slug", slug))
	}

	resp := ErrorResponse{Slug: slug, HttpStatus: status}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
		if logger != nil {
			logger.Error("Failed to encode error response", encodeErr)
		}
	}
}
