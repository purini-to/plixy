package middleware

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/purini-to/plixy/pkg/log"
)

func Recover(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				logger := log.FromContext(r.Context())
				if logger == nil {
					logger = log.GetLogger()
				}

				err, ok := rvr.(error)
				if !ok {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					logger.Error("Unable to handle due to unknown error", zap.Any("unknownErr", rvr))
					return
				}

				http.Error(w, err.Error(), http.StatusInternalServerError)
				logger.Error("Recovers because an error has occurred", zap.Error(err))
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
