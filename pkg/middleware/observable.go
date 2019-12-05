package middleware

import (
	"net/http"

	"github.com/purini-to/plixy/pkg/config"
	ptrace "github.com/purini-to/plixy/pkg/trace"
	"go.opencensus.io/trace"

	"github.com/purini-to/plixy/pkg/stats"

	"go.opencensus.io/tag"

	"go.opencensus.io/plugin/ochttp"
)

func Observable(next http.Handler) http.Handler {
	if config.Global.Trace.Enable {
		next = spanWithRequestID(next)
	}
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := tag.New(r.Context(), tag.Upsert(stats.KeyPath, r.URL.Path))

		handler := &ochttp.Handler{Handler: next}
		handler.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func spanWithRequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		span := trace.FromContext(ctx)
		if span == nil {
			next.ServeHTTP(w, r)
			return
		}

		reqID := ptrace.RequestIDFromContext(ctx)
		if reqID == "" {
			next.ServeHTTP(w, r)
			return
		}

		span.AddAttributes(trace.StringAttribute("plixy.request_id", reqID))
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
