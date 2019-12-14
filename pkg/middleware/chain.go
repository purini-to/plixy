package middleware

import "net/http"

func Chain(next http.Handler, mw []func(http.Handler) http.Handler) http.Handler {
	l := len(mw) - 1
	for i := range mw {
		next = mw[l-i](next)
	}
	return next
}
