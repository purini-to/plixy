package proxy

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/purini-to/plixy/pkg/config"
	"go.opencensus.io/plugin/ochttp"
	"golang.org/x/net/http2"

	"github.com/purini-to/plixy/pkg/api"

	"github.com/purini-to/plixy/pkg/httperr"

	"github.com/pkg/errors"

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
)

type Middleware func(http.Handler) http.Handler

type Router struct {
	middlewares []Middleware
	proxy       *httputil.ReverseProxy
	server      http.Handler
}

func (r *Router) Use(middlewares ...Middleware) {
	r.middlewares = append(r.middlewares, middlewares...)
	r.server = r.chain(r.proxy)
}

func (r *Router) SetMiddlewares(middlewares ...Middleware) {
	r.middlewares = middlewares
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

func New() (*Router, error) {
	tr := &http.Transport{
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
	if err := http2.ConfigureTransport(tr); err != nil {
		return nil, errors.Wrap(err, "could not create http2 transport")
	}

	var transport http.RoundTripper
	transport = tr
	if config.Global.IsObservable() {
		transport = &ochttp.Transport{Base: tr}
	}

	router := &Router{
		middlewares: make([]Middleware, 0),
		proxy: &httputil.ReverseProxy{
			Director:  createDirector(),
			Transport: transport,
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
				// client canceled
				if err.Error() == "context canceled" {
					httperr.ClientClosedRequest(w, err)
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
				httperr.BadGateway(w)
			},
		},
	}

	router.server = router.chain(router.proxy)
	return router, nil
}

func createDirector() func(r *http.Request) {
	return func(r *http.Request) {
		originalURI := r.RequestURI

		apiDef := api.FromContext(r.Context())
		target := apiDef.Proxy.Upstream.Target
		uri, err := url.Parse(target)
		if err != nil {
			panic(errors.New(fmt.Sprintf("Could not parse upstream uri. uri: %s", target)))
		}

		r.URL.Scheme = uri.Scheme
		r.URL.Host = uri.Host
		r.Host = uri.Host

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
