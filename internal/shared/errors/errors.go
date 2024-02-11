package errors

import "net/http"

// ErrorType represents a specific type of error with a corresponding HTTP status code
type ErrorType struct {
	T        string // Consider if you really need to export this field
	HTTPCode int
}

// Define error types with corresponding HTTP status codes
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

// SlugError represents a detailed error with a type and a slug for more context
type SlugError struct {
	error     string
	slug      string
	errorType ErrorType
}

func (s SlugError) Error() string {
	return s.error
}

func (s SlugError) Slug() string {
	return s.slug
}

func (s SlugError) ErrorType() ErrorType {
	return s.errorType
}

// HTTPStatus returns the HTTP status code associated with the SlugError's error type
func (s SlugError) HTTPStatus() int {
	return s.errorType.HTTPCode
}

// Factory functions for creating specific types of SlugErrors
func NewSlugError(errorMsg, slug string, errorType ErrorType) SlugError {
	return SlugError{
		error:     errorMsg,
		slug:      slug,
		errorType: errorType,
	}
}

// Specific error constructors
func NewAuthorizationError(errorMsg, slug string) SlugError {
	return NewSlugError(errorMsg, slug, ErrorTypeAuthorization)
}

func NewNotFoundError(errorMsg, slug string) SlugError {
	return NewSlugError(errorMsg, slug, ErrorTypeNotFound)
}

func NewIncorrectInputError(errorMsg, slug string) SlugError {
	return NewSlugError(errorMsg, slug, ErrorTypeIncorrectInput)
}

func NewDatabaseError(errorMsg, slug string) SlugError {
	return NewSlugError(errorMsg, slug, ErrorTypeDatabase)
}

func NewTimeoutError(errorMsg, slug string) SlugError {
	return NewSlugError(errorMsg, slug, ErrorTypeTimeout)
}

func NewConflictError(errorMsg, slug string) SlugError {
	return NewSlugError(errorMsg, slug, ErrorTypeConflict)
}

func NewRateLimitError(errorMsg, slug string) SlugError {
	return NewSlugError(errorMsg, slug, ErrorTypeRateLimit)
}
