package proxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/purini-to/plixy/pkg/log"
	"go.uber.org/zap"
)

const (
	// DefaultDialTimeout when connecting to a backend server.
	DefaultDialTimeout = 30 * time.Second

	// DefaultIdleConnsPerHost the default value set for http.Transport.MaxIdleConnsPerHost.
	DefaultIdleConnsPerHost = 64

	// DefaultIdleConnTimeout is the default value for the the maximum amount of time an idle
	// (keep-alive) connection will remain idle before closing itself.
	DefaultIdleConnTimeout = 90 * time.Second

	// HTTPStatusClientClosedRequest is status for client is closed
	HTTPStatusClientClosedRequest = 499
)

type middleware func(http.Handler) http.Handler

type Router struct {
	middlewares []middleware
	proxy       *httputil.ReverseProxy
	server      http.Handler
}

func (r *Router) Use(middlewares ...middleware) {
	r.middlewares = append(r.middlewares, middlewares...)
	r.server = r.chain(r.proxy)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.server.ServeHTTP(w, req)
}

func (r *Router) chain(handle http.Handler) http.Handler {
	l := len(r.middlewares) - 1
	for i := range r.middlewares {
		handle = r.middlewares[l-i](handle)
	}
	return handle
}

func New() *Router {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   DefaultDialTimeout,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          256,
		IdleConnTimeout:       DefaultIdleConnTimeout,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: DefaultDialTimeout,
		MaxIdleConnsPerHost:   DefaultIdleConnsPerHost,
	}

	router := &Router{
		middlewares: make([]middleware, 0),
		proxy: &httputil.ReverseProxy{
			Director:  createDirector(),
			Transport: transport,
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
				// client canceled
				if err.Error() == "context canceled" {
					http.Error(w, err.Error(), HTTPStatusClientClosedRequest)
					return
				}

				logger := log.FromContext(r.Context())
				// disabled stacktrace
				logger.WithOptions(zap.AddStacktrace(zapcore.PanicLevel)).
					Error("Error proxy response",
						zap.String("method", r.Method),
						zap.String("upstream_host", r.URL.Host),
						zap.String("upstream_uri", r.RequestURI),
						zap.String("upstream_scheme", r.URL.Scheme),
						zap.Error(err),
					)
				http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
			},
		},
	}
	router.server = router.proxy
	return router
}

func createDirector() func(r *http.Request) {
	return func(r *http.Request) {
		originalURI := r.RequestURI

		r.URL.Scheme = "http"
		r.URL.Host = "dummy.restapiexample.com"
		r.Host = "dummy.restapiexample.com"

		logger := log.FromContext(r.Context())
		logger.Info("Proxying request to the following upstream",
			zap.String("uri", originalURI),
			zap.String("method", r.Method),
			zap.String("upstream_host", r.URL.Host),
			zap.String("upstream_uri", r.URL.RequestURI()),
			zap.String("upstream_scheme", r.URL.Scheme),
		)
	}
}
