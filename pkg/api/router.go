package api

import (
	"context"
	"net/http"
	"sync"

	"github.com/purini-to/plixy/pkg/httperr"
	"github.com/purini-to/plixy/pkg/log"
	"go.opencensus.io/tag"
	"go.opencensus.io/trace"
	"go.uber.org/zap"

	"github.com/gorilla/mux"
	"github.com/purini-to/plixy/pkg/config"
	pstats "github.com/purini-to/plixy/pkg/stats"
	"go.opencensus.io/stats"
)

type Router struct {
	apiConfigMap sync.Map
	mux          *mux.Router
}

func (r *Router) WithApiDefinition(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		var match mux.RouteMatch
		if ok := r.mux.Match(req, &match); !ok || match.MatchErr != nil {
			httperr.Error(w, req, match.MatchErr)
			return
		}

		v, ok := r.apiConfigMap.Load(match.Route.GetName())
		if !ok {
			httperr.NotFound(w)
			return
		}
		apiDef := v.(*Api)

		ctx := ToContext(req.Context(), apiDef)
		log.FromContext(ctx).Debug("Match proxy api", zap.String("name", apiDef.Name))
		if config.Global.Stats.Enable {
			ctx, _ = tag.New(ctx, tag.Upsert(pstats.KeyApiName, apiDef.Name))
		}
		if config.Global.Trace.Enable {
			if span := trace.FromContext(ctx); span != nil {
				span.AddAttributes(trace.StringAttribute("plixy.api_name", apiDef.Name))
			}
		}

		req = req.WithContext(ctx)
		next.ServeHTTP(w, req.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func NewRouter(def *Definition) *Router {
	r := &Router{
		apiConfigMap: sync.Map{},
	}

	m := mux.NewRouter()
	for _, a := range def.Apis {
		m.Name(a.Name).Methods(a.Proxy.Methods...).Path(a.Proxy.Path)
		r.apiConfigMap.Store(a.Name, a)
	}
	r.mux = m

	if config.Global.Stats.Enable {
		stats.Record(context.Background(), pstats.ApiDefinitionVersion.M(def.Version))
	}
	return r
}
