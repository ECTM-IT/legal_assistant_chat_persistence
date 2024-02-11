package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Assume these constants are defined elsewhere in the package.

func decodeJSON(w http.ResponseWriter, r *http.Request, dst interface{}, disallowUnknownFields bool) error {
	// Ensure the request's Content-Type is appropriate
	if r.Header.Get("Content-Type") != "application/json" {
		return errors.New("incorrect content type")
	}

	// Set max bytes for the reader
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	if disallowUnknownFields {
		dec.DisallowUnknownFields()
	}

	if err := dec.Decode(dst); err != nil {
		return handleError(err)
	}

	// Check for additional data after the decoded object
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

// Extracted error handling for readability
func handleError(err error) error {
	maxBytes := 1_048_576
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError

	switch {
	case errors.As(err, &syntaxError):
		return fmt.Errorf("badly-formed JSON: %w", err)
	case errors.As(err, &unmarshalTypeError):
		return fmt.Errorf("incorrect JSON type: %w", err)
	case errors.Is(err, io.ErrUnexpectedEOF), errors.Is(err, io.EOF):
		return errors.New("incomplete or empty JSON")
	case strings.HasPrefix(err.Error(), "json: unknown field"):
		fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
		return fmt.Errorf("unknown field: %s", fieldName)
	case err.Error() == "http: request body too large":
		return fmt.Errorf("body too large, limit is %d bytes", maxBytes)
	default:
		return err // Unknown error
	}
}
