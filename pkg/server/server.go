package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/purini-to/plixy/pkg/api"

	"github.com/purini-to/plixy/pkg/config"

	"github.com/purini-to/plixy/pkg/middleware"
	"github.com/purini-to/plixy/pkg/proxy"

	"github.com/pkg/errors"
	"github.com/purini-to/plixy/pkg/log"
	"go.uber.org/zap"
)

var defaultGraceTimeOut = time.Second * 30

type Server struct {
	server   *http.Server
	proxy    *proxy.Router
	stopChan chan struct{}
	defChan  chan *api.DefinitionChanged
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		defer s.Stop()
		<-ctx.Done()
		log.Info("Stopping server gracefully")
	}()

	middlewares, err := s.buildMiddlewares()
	if err != nil {
		return errors.Wrap(err, "could not build proxy middlewares")
	}

	s.proxy, err = proxy.New()
	if err != nil {
		return errors.Wrap(err, "error proxy.New()")
	}

	s.proxy.Use(middlewares...)

	address := fmt.Sprintf(":%v", config.Global.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrap(err, "error opening listener")
	}

	s.server = &http.Server{
		Handler: s.proxy,
	}

	go func() {
		if err := s.serve(listener); err != http.ErrServerClosed {
			log.Fatal("Could not start http server", zap.Error(err))
		}
	}()
	go s.listenApiDefinition(ctx)

	api.Watch(ctx, s.defChan)

	log.Info("Listening HTTP server", zap.String("address", address))
	return nil
}

func (s *Server) Stop() {
	defer log.Info("Server stopped")

	graceTimeOut, err := time.ParseDuration(config.Global.GraceTimeOut)
	if err != nil {
		log.Error(fmt.Sprintf("Could not parse duration for %s. Set default duration %v",
			config.Global.GraceTimeOut, defaultGraceTimeOut))
		graceTimeOut = defaultGraceTimeOut
	}
	ctx, cancel := context.WithTimeout(context.Background(), graceTimeOut)
	defer cancel()
	log.Info(fmt.Sprintf("Waiting %s before killing connections...", graceTimeOut))
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

func (s *Server) serve(listener net.Listener) error {
	return s.server.Serve(listener)
}

func (s *Server) buildMiddlewares() ([]proxy.Middleware, error) {
	middlewares := []proxy.Middleware{
		middleware.WithLogger(log.GetLogger()),
		middleware.RequestID,
		middleware.AccessLog,
		middleware.Recover,
	}
	if config.Global.Debug {
		middlewares = append(middlewares, middleware.ProxyStats)
	}
	withApiConfig, err := middleware.WithApiConfig()
	if err != nil {
		return nil, errors.Wrap(err, "error middleware.WithApiConfig()")
	}
	middlewares = append(middlewares, withApiConfig)

	return middlewares, nil
}

func (s *Server) listenApiDefinition(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case def, ok := <-s.defChan:
			if !ok {
				return
			}

			s.handleApiDefinitionEvent(def.Definition)
		}
	}
}

func (s *Server) handleApiDefinitionEvent(def *api.Definition) {
	middlewares, err := s.buildMiddlewares()
	if err != nil {
		log.Error("Failed handle event at change api definition")
		return
	}
	s.proxy.SetMiddlewares(middlewares...)
	log.Info("Reloaded proxy based on new api definition")
}

func New() *Server {
	return &Server{
		stopChan: make(chan struct{}),
		defChan:  make(chan *api.DefinitionChanged),
	}
}
