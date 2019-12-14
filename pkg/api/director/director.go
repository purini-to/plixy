package director

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/purini-to/plixy/pkg/api"

	"github.com/purini-to/plixy/pkg/log"
	"go.uber.org/zap"
)

func Director(r *http.Request) {
	ctx := r.Context()
	originalURI := r.RequestURI
	originalPath := r.URL.Path

	apiDef := api.FromContext(ctx)
	target := apiDef.Proxy.Upstream.Target
	uri, err := url.Parse(target)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Could not parse upstream uri. uri: %s", target)))
	}

	r.URL.Scheme = uri.Scheme
	r.URL.Host = uri.Host
	r.Host = uri.Host

	path := uri.Path
	if len(apiDef.Proxy.Upstream.Vars) > 0 {
		vars := api.VarsFromContext(ctx)
		for _, v := range apiDef.Proxy.Upstream.Vars {
			if s, ok := vars[v]; ok {
				path = strings.Replace(path, fmt.Sprintf("{%s}", v), s, 1)
			}
		}
	}

	r.URL.Path = path
	if !apiDef.Proxy.Upstream.FixedPath {
		r.URL.Path += originalPath
	}
	r.URL.Path = strings.ReplaceAll(r.URL.Path, "//", "/")

	logger := log.FromContext(ctx)
	logger.Info("Proxying request to the following upstream",
		zap.String("uri", originalURI),
		zap.String("method", r.Method),
		zap.String("upstream_host", r.URL.Host),
		zap.String("upstream_uri", r.URL.RequestURI()),
		zap.String("upstream_scheme", r.URL.Scheme),
	)
}
