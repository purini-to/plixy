package middleware

import (
	"net/http"

	"github.com/purini-to/plixy/pkg/log"

	"go.uber.org/zap"
)

// WithLogger sets the logger to the context
func WithLogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := log.ToContext(r.Context(), logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
