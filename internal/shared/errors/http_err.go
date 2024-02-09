package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/errors"
	"github.com/ECTM-IT/legal_assistant_chat_persistence/internal/shared/logs"
)

// ErrorResponder defines an interface for error handling
type ErrorResponder interface {
	RespondWithError(ctx context.Context, err error, slug string, w http.ResponseWriter, r *http.Request)
}

// Logger allows you to inject different logging implementations
type Logger interface {
	WithError(err error) *logs.Entry
	WithField(key string, value interface{}) *logs.Entry
	Warn(msg string)
}

// DefaultErrorResponder is the production implementation
type DefaultErrorResponder struct {
	Logger Logger
}

// RespondWithError provides response handling logic
func (d *DefaultErrorResponder) RespondWithError(ctx context.Context, err error, slug string, w http.ResponseWriter, r *http.Request) {
	// Enhanced error handling for SlugErrors
	if slugError, ok := err.(errors.SlugError); ok {
		d.handleSlugError(ctx, slugError, w, r)
		return
	}

	// Default to internal server error
	d.internalServerError(ctx, err, slug, w, r)
}

func (d *DefaultErrorResponder) handleSlugError(ctx context.Context, slugError errors.SlugError, w http.ResponseWriter, r *http.Request) {
	// Extract error status directly from SlugError
	status := slugError.HTTPStatus()

	d.httpRespondWithError(ctx, slugError, slugError.Slug(), w, r, slugError.Error(), status)
}

// Helper functions (unchanged, but now methods of DefaultErrorResponder)
func (d *DefaultErrorResponder) internalServerError(ctx context.Context, err error, slug string, w http.ResponseWriter, r *http.Request) {
	d.httpRespondWithError(ctx, err, slug, w, r, "Internal server error", http.StatusInternalServerError)
}

func (d *DefaultErrorResponder) unauthorised(ctx context.Context, err error, slug string, w http.ResponseWriter, r *http.Request) {
	d.httpRespondWithError(ctx, err, slug, w, r, "Unauthorised", http.StatusUnauthorized)
}

// ... similar helper functions for badRequest, notFound

func (d *DefaultErrorResponder) httpRespondWithError(ctx context.Context, err error, slug string, w http.ResponseWriter, r *http.Request, logMsg string, status int) {
	// Use request-scoped logger from the context
	logger := logs.GetLogEntry(r)
	if ctxLogger, ok := logs.FromContext(ctx); ok {
		logger = ctxLogger
	}

	logger.WithError(err).WithField("error-slug", slug).Warn(logMsg)

	resp := ErrorResponse{slug, status}
	w.Header().Set("Content-Type", "application/json") // Explicitly set Content-Type
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// Improved panic handling
		logger.WithError(err).Error("Failed to encode error response")
	}
}

// ErrorResponse (unchanged)
type ErrorResponse struct {
	Slug       string `json:"slug"`
	httpStatus int
}

func (e ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(e.httpStatus)
	return nil
}
