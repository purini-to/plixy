package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/purini-to/plixy/pkg/middleware"

	"github.com/purini-to/plixy/pkg/stats"

	"github.com/pkg/errors"
	"github.com/purini-to/plixy/pkg/config"
	"github.com/purini-to/plixy/pkg/log"
	"go.uber.org/zap"
)

type Exporter struct {
	server *http.Server
}

func (s *Exporter) StartWithContext(ctx context.Context) error {
	go func() {
		defer s.Close()
		<-ctx.Done()
	}()

	address := fmt.Sprintf(":%v", config.Global.Stats.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrap(err, "error opening listener")
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", middleware.Recover(stats.PrometheusExporter))
	s.server = &http.Server{
		Handler: mux,
	}

	go func() {
		if err := s.server.Serve(listener); err != http.ErrServerClosed {
			log.Fatal("Could not start http server", zap.Error(err))
		}
	}()

	log.Info("Listening stats exporter server", zap.String("address", address))
	return nil
}

func (s *Exporter) Close() error {
	return s.server.Close()
}

func NewStatsExporter() *Exporter {
	return &Exporter{}
}
