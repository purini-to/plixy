package middleware

import (
	"net/http"

	"github.com/rs/xid"

	"github.com/purini-to/plixy/pkg/trace"
)

func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		id := trace.RequestIDFromRequest(r)
		if id == "" {
			id = xid.New().String()
		}

		trace.RequestIDToReqRes(w, r, id)
		ctx := trace.RequestIDToContext(r.Context(), id)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
