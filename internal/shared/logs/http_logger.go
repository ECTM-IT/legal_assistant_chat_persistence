package logs

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type contextKey string

const logEntryCtxKey = contextKey("logEntry")

// NewHttpLogger creates a middleware using the Zap logger
func NewHttpLogger(logger *ZapLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a response writer wrapper to capture the status code
			wrappedWriter := newResponseWriter(w)

			// Initialize fields for structured logging (reused later)
			fields := []zap.Field{
				zap.String("method", r.Method),
				zap.String("uri", r.RequestURI),
				zap.String("remote_addr", r.RemoteAddr),
			}

			// Store the logger with initial fields in the request context
			ctx := context.WithValue(r.Context(), logEntryCtxKey, logger.logger.With(fields...))
			r = r.WithContext(ctx)

			// Process the request
			next.ServeHTTP(wrappedWriter, r)

			// Update final status and duration before logging
			fields = append(fields,
				zap.Int("status", wrappedWriter.status),
				zap.Duration("duration", time.Since(start)))

			// Retrieve the logger from the context and log
			GetLogEntry(r).Info("Completed handling request", fields...)
		})
	}
}

// GetLogEntry extracts the Zap logger from the request context
func GetLogEntry(r *http.Request) *zap.Logger {
	entry, ok := r.Context().Value(logEntryCtxKey).(*zap.Logger)
	if !ok {
		return zap.NewNop() // Consider proper error handling or flexibility here
	}
	return entry
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
