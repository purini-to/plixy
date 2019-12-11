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

type Proxy struct {
	server *httputil.ReverseProxy
}

func (r *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.server.ServeHTTP(w, req)
}

func New() (*Proxy, error) {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   config.Global.DialTimeout,
			KeepAlive: config.Global.IdleConnTimeout,
		}).DialContext,
		MaxIdleConns:          config.Global.MaxIdleConns,
		IdleConnTimeout:       config.Global.IdleConnTimeout,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		MaxIdleConnsPerHost:   config.Global.MaxIdleConnsPerHost,
	}
	if err := http2.ConfigureTransport(tr); err != nil {
		return nil, errors.Wrap(err, "could not create http2 transport")
	}

	var transport http.RoundTripper
	transport = tr
	if config.Global.IsObservable() {
		transport = &ochttp.Transport{Base: tr}
	}

	proxy := &Proxy{
		server: &httputil.ReverseProxy{
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

	return proxy, nil
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
