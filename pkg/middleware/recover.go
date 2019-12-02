package middleware

import (
	"net/http"

	"github.com/purini-to/plixy/pkg/httperr"

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
					httperr.InternalServerError(w, http.StatusText(http.StatusInternalServerError))
					logger.Error("Unable to handle due to unknown rvr", zap.Any("rvr", rvr))
					return
				}

				httperr.InternalServerError(w, err.Error())
				logger.Error("Recovers because an error has occurred", zap.Error(err))
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
