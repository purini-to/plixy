package middleware

import (
	"net/http"

	"github.com/purini-to/plixy/pkg/log"
	"go.uber.org/zap"

	"github.com/rs/xid"

	"github.com/purini-to/plixy/pkg/trace"
)

// RequestID set a unique ID in the header
// If there is already an ID in the request header,
// it will be relayed to the header without generating a new one
func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		id := trace.RequestIDFromRequest(r)
		if id == "" {
			id = xid.New().String()
		}

		ctx := r.Context()
		if logger := log.FromContext(ctx); logger != nil {
			logger := logger.With(zap.String("request_id", id))
			ctx = log.ToContext(ctx, logger)
		}

		trace.RequestIDToReqRes(w, r, id)
		ctx = trace.RequestIDToContext(ctx, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
