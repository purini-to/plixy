package middleware

import (
	"net/http"

	"github.com/purini-to/plixy/pkg/stats"

	"go.opencensus.io/tag"

	"go.opencensus.io/plugin/ochttp"
)

func Stats(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := tag.New(r.Context(), tag.Upsert(stats.KeyPath, r.URL.Path))

		handler := &ochttp.Handler{Handler: next}
		handler.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
