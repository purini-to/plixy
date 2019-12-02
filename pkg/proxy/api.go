package proxy

import (
	"net/http"
	"sync"

	"github.com/purini-to/plixy/pkg/httperr"

	"github.com/purini-to/plixy/pkg/config"

	"go.uber.org/zap"

	"github.com/purini-to/plixy/pkg/log"

	"github.com/gorilla/mux"

	"github.com/pkg/errors"
	"github.com/purini-to/plixy/pkg/api"
)

var apiMap = sync.Map{}

func ConfigHandleCreator() (Middleware, error) {
	apis, err := api.GetApiConfigs()
	if err != nil {
		return nil, errors.Wrap(err, "could not api.GetApiConfigs()")
	}

	router := mux.NewRouter()
	for _, a := range apis {
		router.Name(a.Name).Methods(a.Proxy.Methods...).Path(a.Proxy.Path)
		apiMap.Store(a.Name, a)
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var match mux.RouteMatch
			if ok := router.Match(r, &match); !ok || match.MatchErr != nil {
				httperr.Error(w, r, match.MatchErr)
				return
			}

			v, ok := apiMap.Load(match.Route.GetName())
			if !ok {
				httperr.NotFound(w)
				return
			}
			configApi := v.(*config.Api)

			ctx := config.ApiToContext(r.Context(), configApi)

			log.FromContext(ctx).Debug("Match proxy api", zap.String("name", configApi.Name))

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}, nil
}
