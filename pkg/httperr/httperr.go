package httperr

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/purini-to/plixy/pkg/log"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// HTTPStatusClientClosedRequest is status for client is closed
var HTTPStatusClientClosedRequest = 499

func NotFound(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}

func ClientClosedRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), HTTPStatusClientClosedRequest)
}

func BadGateway(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
}

func MethodNotAllowed(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}

func InternalServerError(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusInternalServerError)
}

func Error(w http.ResponseWriter, r *http.Request, err error) {
	cause := errors.Cause(err)
	switch cause.(type) {
	default:
		errorByString(w, r, err, cause)
		return
	}
}

func errorByString(w http.ResponseWriter, r *http.Request, err, cause error) {
	switch cause.Error() {
	case mux.ErrNotFound.Error():
		NotFound(w)
		return
	case mux.ErrMethodMismatch.Error():
		MethodNotAllowed(w)
		return
	case "context canceled":
		ClientClosedRequest(w, cause)
		return
	default:
		InternalServerError(w, cause.Error())
		log.FromContext(r.Context()).Error("Unable to handle due to unknown error", zap.Error(err))
		return
	}
}
