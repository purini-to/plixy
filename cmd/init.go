package cmd

import (
	"time"

	"github.com/purini-to/plixy/pkg/stats"

	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/pkg/errors"
	"github.com/purini-to/plixy/pkg/config"
	"github.com/purini-to/plixy/pkg/log"
	"go.opencensus.io/stats/view"
	"go.uber.org/zap"
)

func initLog() error {
	var (
		logger *zap.Logger
		err    error
	)
	if debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		return errors.Wrap(err, "could not build logger")
	}

	log.SetLogger(logger)
	return nil
}

func initConfig(configFilePath string) error {
	opts := []config.Option{
		config.WithEnvPrefix("PLIXY"),
	}
	if len(configFilePath) > 0 {
		opts = append(opts, config.WithLoadFile(configFilePath))
	}

	err := config.Load(opts...)
	if err != nil {
		return errors.Wrap(err, "error load config")
	}
	log.Debug("Load config", zap.Any("config", config.Global))
	return nil
}

func initExporter() error {
	if config.Global.Stats.Enable {
		if err := initStatsExporter(); err != nil {
			return errors.Wrap(err, "failed to create stats exporter")
		}
	}

	return nil
}

func initStatsExporter() error {
	exporter, err := prometheus.NewExporter(prometheus.Options{
		Namespace: config.Global.Stats.ServiceName,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create prometheus exporter")
	}

	view.RegisterExporter(exporter)
	stats.PrometheusExporter = exporter

	view.SetReportingPeriod(5 * time.Second)

	if err := view.Register(stats.AllViews...); err != nil {
		return errors.Wrap(err, "failed to register server views")
	}

	return nil
}
