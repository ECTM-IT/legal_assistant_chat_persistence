package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// APIError represents a standardized error response.
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// BaseHandler provides utility methods for HTTP handlers,
// including JSON response writing and input parsing.
type BaseHandler struct{}

// RespondWithJSON writes the given payload to the ResponseWriter as JSON.
// It sets the Content-Type to "application/json" and writes the HTTP status code.
func (h *BaseHandler) RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		h.RespondWithError(w, http.StatusInternalServerError, "Failed to marshal JSON response")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(response)
}

// RespondWithError writes an error response in a standardized JSON structure.
// It includes an HTTP code and a descriptive error message.
func (h *BaseHandler) RespondWithError(w http.ResponseWriter, code int, message string) {
	errorPayload := APIError{
		Code:    code,
		Message: message,
	}
	h.RespondWithJSON(w, code, errorPayload)
}

// ParseObjectID extracts a MongoDB ObjectID from either the request headers or URL variables.
// If fromHeader is true, it attempts to parse the "Authorization" header as an ObjectID.
// Otherwise, it attempts to parse the given key in the URL variables.
// Returns an error if the ObjectID is invalid or missing.
func (h *BaseHandler) ParseObjectID(r *http.Request, key string, fromHeader bool) (primitive.ObjectID, error) {
	var idStr string
	if fromHeader {
		idStr = r.Header.Get("Authorization")
	} else {
		idStr = mux.Vars(r)[key]
	}

	if idStr == "" {
		return primitive.NilObjectID, errors.New("object ID string is empty")
	}

	objectID, err := primitive.ObjectIDFromHex(strings.TrimSpace(idStr))
	if err != nil {
		return primitive.NilObjectID, errors.New("invalid ObjectID provided")
	}

	return objectID, nil
}

// ParseURLVar retrieves the value of the specified URL variable, trims whitespace,
// and returns it as a string. If the variable does not exist or is empty, it returns an empty string.
// Note: Consider validating the returned value as needed in the handler logic.
func (h *BaseHandler) ParseURLVar(r *http.Request, key string) string {
	return strings.TrimSpace(mux.Vars(r)[key])
}

// DecodeJSONBody decodes the JSON body of the request into the provided interface.
// It returns an error if the JSON is malformed or if decoding fails.
// The caller is responsible for handling the error and ensuring the request body is closed.
func (h *BaseHandler) DecodeJSONBody(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Encourage strict decoding to avoid silent failures.
	if err := decoder.Decode(v); err != nil {
		return errors.New("invalid JSON payload")
	}
	return nil
}
