package errors

import "net/http"

// ErrorType represents a specific type of error with a corresponding HTTP status code.
type ErrorType struct {
	Type     string
	HTTPCode int
}

// Define error types with corresponding HTTP status codes.
var (
	ErrorTypeUnknown        = ErrorType{"unknown", http.StatusInternalServerError}
	ErrorTypeAuthorization  = ErrorType{"authorization", http.StatusUnauthorized}
	ErrorTypeIncorrectInput = ErrorType{"incorrect-input", http.StatusBadRequest}
	ErrorTypeNotFound       = ErrorType{"not-found", http.StatusNotFound}
	ErrorTypeDatabase       = ErrorType{"database", http.StatusInternalServerError}
	ErrorTypeTimeout        = ErrorType{"timeout", http.StatusGatewayTimeout}
	ErrorTypeConflict       = ErrorType{"conflict", http.StatusConflict}
	ErrorTypeRateLimit      = ErrorType{"rate-limit", http.StatusTooManyRequests}
)

// SlugError represents a detailed error with a type and a slug for more context.
type SlugError struct {
	message   string
	slug      string
	errorType ErrorType
}

// Error implements the error interface for SlugError.
func (s SlugError) Error() string {
	return s.message
}

// Slug returns the slug associated with the SlugError.
func (s SlugError) Slug() string {
	return s.slug
}

// ErrorType returns the error type of the SlugError.
func (s SlugError) ErrorType() ErrorType {
	return s.errorType
}

// HTTPStatus returns the HTTP status code associated with the SlugError's error type.
func (s SlugError) HTTPStatus() int {
	return s.errorType.HTTPCode
}

// NewSlugError creates a new SlugError with the given message, slug, and error type.
func NewSlugError(message, slug string, errorType ErrorType) SlugError {
	return SlugError{
		message:   message,
		slug:      slug,
		errorType: errorType,
	}
}

// Specific error constructors for creating different types of SlugErrors.
func NewAuthorizationError(message, slug string) SlugError {
	return NewSlugError(message, slug, ErrorTypeAuthorization)
}

func NewNotFoundError(message, slug string) SlugError {
	return NewSlugError(message, slug, ErrorTypeNotFound)
}

func NewIncorrectInputError(message, slug string) SlugError {
	return NewSlugError(message, slug, ErrorTypeIncorrectInput)
}

func NewDatabaseError(message, slug string) SlugError {
	return NewSlugError(message, slug, ErrorTypeDatabase)
}

func NewTimeoutError(message, slug string) SlugError {
	return NewSlugError(message, slug, ErrorTypeTimeout)
}

func NewConflictError(message, slug string) SlugError {
	return NewSlugError(message, slug, ErrorTypeConflict)
}

func NewRateLimitError(message, slug string) SlugError {
	return NewSlugError(message, slug, ErrorTypeRateLimit)
}
