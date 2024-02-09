package logs

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type contextKey string

const logEntryCtxKey = contextKey("logEntry")

func NewHttpLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a response writer wrapper to capture the status code
			wrappedWriter := newResponseWriter(w)

			// Logging logic
			entry := logger.WithFields(logrus.Fields{
				"status":      wrappedWriter.status,
				"method":      r.Method,
				"uri":         r.RequestURI,
				"remote_addr": r.RemoteAddr,
				"duration":    time.Since(start),
				// "req_id":       can be added if you have a request ID
			})

			ctx := context.WithValue(r.Context(), logEntryCtxKey, entry)
			r = r.WithContext(ctx)

			// Process the request
			next.ServeHTTP(wrappedWriter, r)

			entry.Info("Completed handling request")
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func GetLogEntry(r *http.Request) logrus.FieldLogger {
	entry, ok := r.Context().Value(logEntryCtxKey).(logrus.FieldLogger)
	if !ok {
		// Log entry not found in context, return a new one
		return logrus.NewEntry(logrus.StandardLogger())
	}
	return entry
}
