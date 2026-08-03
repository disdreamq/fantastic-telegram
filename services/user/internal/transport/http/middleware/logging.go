package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
	body   []byte
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}
	rw.body = append(rw.body, b...)
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

func LoggingMiddleware(logger zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			traceID := r.Header.Get("X-Request-ID")
			if traceID == "" {
				traceID = uuid.New().String()
			}

			ctx := context.WithValue(r.Context(), "trace_id", traceID)
			ctx = logger.WithContext(ctx)
			r = r.WithContext(ctx)

			w.Header().Set("X-Request-ID", traceID)

			rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(rw, r)

			duration := time.Since(start)

			if rw.status >= 500 {
				logger.Error().
					Str("trace_id", traceID).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Int("status", rw.status).
					Dur("duration", duration).
					Str("response_body", string(rw.body)).
					Msg("Request failed")
			} else if rw.status >= 400 {
				logger.Warn().
					Str("trace_id", traceID).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Int("status", rw.status).
					Str("response_body", string(rw.body)).
					Msg("Request failed")
			} else {
				logger.Info().
					Str("trace_id", traceID).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Int("status", rw.status).
					Dur("duration", duration).
					Str("remote_addr", r.RemoteAddr).
					Msg("HTTP request")
			}
		})
	}
}
