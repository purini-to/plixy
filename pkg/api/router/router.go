package router

import (
	"net/http"

	"github.com/purini-to/plixy/pkg/middleware"

	"github.com/purini-to/plixy/pkg/plugin"

	"github.com/purini-to/plixy/pkg/api"

	"github.com/purini-to/plixy/pkg/httperr"
	"github.com/purini-to/plixy/pkg/log"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
	"go.uber.org/zap"

	"github.com/gorilla/mux"
	"github.com/purini-to/plixy/pkg/config"
	pstats "github.com/purini-to/plixy/pkg/stats"
)

type Route struct {
	api *api.Api
	mw  []func(next http.Handler) http.Handler
}

type Router struct {
	apiConfigMap map[string]*Route
	mux          *mux.Router
}

func (r *Router) WithApiDefinition(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		var match mux.RouteMatch
		if ok := r.mux.Match(req, &match); !ok || match.MatchErr != nil {
			httperr.Error(w, req, match.MatchErr)
			return
		}

		v, ok := r.apiConfigMap[match.Route.GetName()]
		if !ok {
			httperr.NotFound(w)
			return
		}
		apiDef := v.api

		ctx := api.ToContext(req.Context(), apiDef)
		ctx = api.VarsToContext(ctx, match.Vars)
		log.FromContext(ctx).Debug("Match proxy api", zap.String("name", apiDef.Name))
		if config.Global.Stats.Enable {
			ctx, _ = tag.New(ctx, tag.Upsert(pstats.KeyApiName, apiDef.Name))
		}
		if config.Global.Trace.Enable {
			if span := trace.FromContext(ctx); span != nil {
				span.AddAttributes(trace.StringAttribute("plixy.api_name", apiDef.Name))
			}
		}

		req.Header.Set(api.NameHeaderKey, apiDef.Name)

		req = req.WithContext(ctx)
		//next.ServeHTTP(w, req)
		middleware.Chain(next, v.mw).ServeHTTP(w, req)
	}
	return http.HandlerFunc(fn)
}

func NewRouter(def *api.Definition) (*Router, error) {
	r := &Router{
		apiConfigMap: make(map[string]*Route, 0),
	}

	m := mux.NewRouter()
	for _, a := range def.Apis {
		rt := m.Name(a.Name).Path(a.Proxy.Path)
		if len(a.Proxy.Methods) > 0 {
			rt = rt.Methods(a.Proxy.Methods...)
		}

		handlers, err := plugin.BuildBeforeProxy(a.Plugins)
		if err != nil {
			return nil, err
		}
		r.apiConfigMap[a.Name] = &Route{
			api: a,
			mw:  handlers,
		}
	}
	r.mux = m

	return r, nil
}
