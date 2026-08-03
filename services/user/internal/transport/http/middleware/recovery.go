package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog/log"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				panicMsg := fmt.Sprintf("%v", err)
				if e, ok := err.(error); ok {
					panicMsg = e.Error()
				}
				log.Error().
					Str("panic", panicMsg).
					Str("stack", string(debug.Stack())).
					Msg("panic recovered")

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
