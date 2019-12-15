package plugin

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/pkg/errors"

	"github.com/purini-to/plixy/pkg/api"

	"github.com/purini-to/plixy/pkg/log"
	"go.uber.org/zap"
)

type cache struct {
	validateConfig sync.Map
	beforeProxy    sync.Map
}

var registered = &cache{}

type BeforeProxyFunc func(config map[string]interface{}) (func(next http.Handler) http.Handler, error)

type Plugin struct {
	BeforeProxy BeforeProxyFunc
}

func Register(name string, plg *Plugin) {
	log.Debug("Register plugin", zap.String("name", name))

	if plg.BeforeProxy != nil {
		registered.beforeProxy.Store(name, plg.BeforeProxy)
	}
}

func BuildBeforeProxy(plg []*api.Plugin) ([]func(next http.Handler) http.Handler, error) {
	mw := make([]func(next http.Handler) http.Handler, 0)
	for _, p := range plg {
		value, ok := registered.beforeProxy.Load(p.Name)
		if !ok {
			return nil, errors.New(fmt.Sprintf("not found plugin. name: %s", p.Name))
		}
		f := value.(BeforeProxyFunc)
		h, err := f(p.Config)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed BeforeProxy plugin. name: %s", p.Name))
		}
		mw = append(mw, h)
	}

	return mw, nil
}
