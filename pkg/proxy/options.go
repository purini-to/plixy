package proxy

import (
	"net/http"
)

type Option func(opt *proxyOpt)

func WithTransport(v http.RoundTripper) Option {
	return func(t *proxyOpt) {
		t.transport = v
	}
}

func WithDirector(v func(*http.Request)) Option {
	return func(t *proxyOpt) {
		t.director = v
	}
}

func WithErrorHandler(v func(w http.ResponseWriter, r *http.Request, err error)) Option {
	return func(t *proxyOpt) {
		t.errorHandler = v
	}
}
