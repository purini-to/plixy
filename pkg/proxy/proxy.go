package proxy

import (
	"net/http"
	"net/http/httputil"

	"github.com/purini-to/plixy/pkg/proxy/transport"

	"github.com/purini-to/plixy/pkg/httperr"

	"go.uber.org/zap/zapcore"

	"github.com/purini-to/plixy/pkg/log"
	"go.uber.org/zap"
)

var (
	defaultErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
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
	}
)

type Proxy struct {
	server *httputil.ReverseProxy
}

func (r *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.server.ServeHTTP(w, req)
}

type proxyOpt struct {
	transport    http.RoundTripper
	director     func(*http.Request)
	errorHandler func(w http.ResponseWriter, r *http.Request, err error)
}

func New(opts ...Option) (*Proxy, error) {
	p := proxyOpt{}

	for _, opt := range opts {
		opt(&p)
	}

	if p.transport == nil {
		p.transport = transport.New()
	}

	if p.errorHandler == nil {
		p.errorHandler = defaultErrorHandler
	}

	proxy := &Proxy{
		server: &httputil.ReverseProxy{
			Director:     p.director,
			Transport:    p.transport,
			ErrorHandler: p.errorHandler,
		},
	}

	return proxy, nil
}
