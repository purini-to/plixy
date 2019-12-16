package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/purini-to/plixy/pkg/store"

	"github.com/purini-to/plixy/pkg/api/router"

	"github.com/purini-to/plixy/pkg/health"

	"github.com/purini-to/plixy/pkg/api"

	"github.com/purini-to/plixy/pkg/config"

	"github.com/purini-to/plixy/pkg/middleware"
	"github.com/purini-to/plixy/pkg/proxy"

	"github.com/pkg/errors"
	"github.com/purini-to/plixy/pkg/log"
	"go.uber.org/zap"

	// plugins
	_ "github.com/purini-to/plixy/pkg/plugin/rate"
)

type Server struct {
	sync.RWMutex
	server      *http.Server
	proxy       *proxy.Proxy
	router      *router.Router
	middlewares []func(http.Handler) http.Handler
	store       store.Store
	stopChan    chan struct{}
	defChan     chan *api.DefinitionChanged
}

func New(store store.Store) *Server {
	return &Server{
		store:    store,
		stopChan: make(chan struct{}),
		defChan:  make(chan *api.DefinitionChanged),
	}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		defer s.Stop()
		<-ctx.Done()
		log.Info("Stopping server gracefully")
	}()

	if err := s.buildMiddlewares(); err != nil {
		return errors.Wrap(err, "could not build server middlewares")
	}

	def, err := s.store.GetDefinition()
	if err != nil {
		return errors.Wrap(err, "could not get api definition")
	}
	rt, err := router.NewRouter(def)
	if err != nil {
		return err
	}
	s.router = rt

	s.proxy, err = proxy.New()
	if err != nil {
		return errors.Wrap(err, "error proxy.New()")
	}

	if config.Global.Watch {
		if err = s.store.Watch(ctx, config.Global.WatchInterval, s.defChan); err != nil {
			return errors.Wrap(err, "Could not watch the api definition")
		}
	}

	address := fmt.Sprintf(":%v", config.Global.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrap(err, "error opening listener")
	}

	s.server = &http.Server{
		Handler: s.buildMux(),
	}

	go func() {
		if err := s.serve(listener); err != http.ErrServerClosed {
			log.Fatal("Could not start http server", zap.Error(err))
		}
	}()
	go s.listenApiDefinition(ctx)

	log.Info("Listening HTTP server", zap.String("address", address))
	return nil
}

func (s *Server) Stop() {
	defer log.Info("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), config.Global.GraceTimeOut)
	defer cancel()
	log.Info(fmt.Sprintf("Waiting %s before killing connections...", config.Global.GraceTimeOut))
	if err := s.server.Shutdown(ctx); err != nil {
		log.Debug("Wait is over due to error", zap.Error(err))
		_ = s.server.Close()
	}
	log.Debug("Server closed")

	s.stopChan <- struct{}{}
}

func (s *Server) Close() error {
	defer close(s.stopChan)
	defer close(s.defChan)
	return s.server.Close()
}

func (s *Server) Wait() {
	log.Debug("Start waiting")
	<-s.stopChan
	log.Debug("Waiting has ended")
}

func (s *Server) buildMux() http.Handler {
	next := s.router.WithApiDefinition(s.proxy)
	next = middleware.Chain(next, s.middlewares)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/__health__" {
			health.Handler(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) buildMiddlewares() error {
	mw := []func(http.Handler) http.Handler{
		middleware.WithLogger(log.GetLogger()),
		middleware.RequestID,
		middleware.RealIP,
		middleware.AccessLog,
	}

	if config.Global.Debug {
		mw = append(mw, middleware.ProxyStats)
	}
	if config.Global.IsObservable() {
		mw = append(mw, middleware.Observable)
	}
	mw = append(mw, middleware.Recover)

	s.middlewares = mw
	return nil
}

func (s *Server) serve(listener net.Listener) error {
	return s.server.Serve(listener)
}

func (s *Server) listenApiDefinition(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case _, ok := <-s.defChan:
			if !ok {
				return
			}
			s.handleApiDefinitionEvent()
		}
	}
}

func (s *Server) handleApiDefinitionEvent() {
	s.Lock()
	defer s.Unlock()
	def, err := s.store.GetDefinition()
	if err != nil {
		log.Error("failed get definition", zap.Error(err))
		return
	}
	rt, err := router.NewRouter(def)
	if err != nil {
		log.Error("could not new router", zap.Error(err))
		return
	}
	s.router = rt
	s.server.Handler = s.buildMux()
	log.Info("Reloaded proxy based on new api definition")
}
