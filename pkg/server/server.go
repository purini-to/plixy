package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

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
	stopChan chan struct{}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		defer s.Stop()
		<-ctx.Done()
		log.Info("Stopping server gracefully")
	}()

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
		return errors.Wrap(err, "error middleware.WithApiConfig()")
	}
	middlewares = append(middlewares, withApiConfig)

	r, err := proxy.New()
	if err != nil {
		return errors.Wrap(err, "error proxy.New()")
	}

	r.Use(middlewares...)

	address := fmt.Sprintf(":%v", config.Global.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrap(err, "error opening listener")
	}

	s.server = &http.Server{
		Handler: r,
	}

	go func() {
		if err := s.serve(listener); err != http.ErrServerClosed {
			log.Fatal("Could not start http server", zap.Error(err))
		}
	}()

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

func New() *Server {
	return &Server{
		stopChan: make(chan struct{}, 1),
	}
}
