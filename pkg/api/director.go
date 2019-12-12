package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/purini-to/plixy/pkg/log"
	"go.uber.org/zap"
)

func Director(r *http.Request) {
	originalURI := r.RequestURI
	originalPath := r.URL.Path

	apiDef := FromContext(r.Context())
	target := apiDef.Proxy.Upstream.Target
	uri, err := url.Parse(target)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Could not parse upstream uri. uri: %s", target)))
	}

	r.URL.Scheme = uri.Scheme
	r.URL.Host = uri.Host
	r.Host = uri.Host

	r.URL.Path = uri.Path
	if !apiDef.Proxy.FixedPath {
		r.URL.Path += originalPath
	}

	logger := log.FromContext(r.Context())
	logger.Info("Proxying request to the following upstream",
		zap.String("uri", originalURI),
		zap.String("method", r.Method),
		zap.String("upstream_host", r.URL.Host),
		zap.String("upstream_uri", r.URL.RequestURI()),
		zap.String("upstream_scheme", r.URL.Scheme),
	)
}
