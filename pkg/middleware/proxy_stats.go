package middleware

import (
	"net/http"

	"github.com/purini-to/plixy/pkg/log"
	"go.uber.org/zap"

	httpstat "github.com/tcnksm/go-httpstat"
)

func ProxyStats(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		result := new(httpstat.Result)
		ctx := httpstat.WithHTTPStat(r.Context(), result)
		defer func() {
			logger := log.FromContext(ctx)
			logger.Debug("Proxy stats", zap.Any("result", result))
		}()

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
