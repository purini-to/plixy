package middleware

import (
	"context"
	"net/http"
	"sync"

	"github.com/purini-to/plixy/pkg/config"
	pstats "github.com/purini-to/plixy/pkg/stats"
	"go.opencensus.io/stats"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/purini-to/plixy/pkg/api"
	"github.com/purini-to/plixy/pkg/httperr"
	"github.com/purini-to/plixy/pkg/log"
	"go.uber.org/zap"
)

var apiConfigMap = sync.Map{}

func WithApiConfig() (func(next http.Handler) http.Handler, error) {
	def, err := api.GetDefinition()
	if err != nil {
		return nil, errors.Wrap(err, "could not api.GetDefinition()")
	}
	if config.Global.Stats.Enable {
		stats.Record(context.Background(), pstats.ApiDefinitionVersion.M(def.Version))
	}

	router := mux.NewRouter()
	for _, a := range def.Apis {
		router.Name(a.Name).Methods(a.Proxy.Methods...).Path(a.Proxy.Path)
		apiConfigMap.Store(a.Name, a)
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var match mux.RouteMatch
			if ok := router.Match(r, &match); !ok || match.MatchErr != nil {
				httperr.Error(w, r, match.MatchErr)
				return
			}

			v, ok := apiConfigMap.Load(match.Route.GetName())
			if !ok {
				httperr.NotFound(w)
				return
			}
			configApi := v.(*api.Api)

			ctx := api.ToContext(r.Context(), configApi)

			log.FromContext(ctx).Debug("Match proxy api", zap.String("name", configApi.Name))

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}, nil
}
