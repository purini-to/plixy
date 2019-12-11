package stats

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/pkg/errors"
	"github.com/purini-to/plixy/pkg/config"
	"github.com/purini-to/plixy/pkg/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type PrometheusExporter struct {
	*prometheus.Exporter
	port   uint
	server *http.Server
}

func (p *PrometheusExporter) Start(ctx context.Context) error {
	go func() {
		defer p.Close()
		<-ctx.Done()
	}()

	address := fmt.Sprintf(":%v", config.Global.Stats.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrap(err, "error opening listener for prometheus")
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", p)
	p.server = &http.Server{
		Handler: mux,
	}

	go func() {
		if err := p.server.Serve(listener); err != http.ErrServerClosed {
			log.Fatal("Could not start prometheus metrics server", zap.Error(err))
		}
	}()

	log.Info("Listening prometheus stats exporter server", zap.String("address", address))
	return nil
}

func (p *PrometheusExporter) Close() error {
	return p.server.Close()
}

type PrometheusOption struct {
	Namespace string
	Port      uint
}

func NewPrometheusExporter(opt *PrometheusOption) (*PrometheusExporter, error) {
	exporter, err := prometheus.NewExporter(prometheus.Options{
		Namespace: opt.Namespace,
		OnError: func(err error) {
			log.GetLogger().WithOptions(zap.AddStacktrace(zapcore.PanicLevel)).Error("Error prometheus exporter", zap.Error(err))
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create prometheus exporter")
	}

	return &PrometheusExporter{
		Exporter: exporter,
		port:     opt.Port,
	}, nil
}
