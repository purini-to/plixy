package cmd

import (
	"github.com/purini-to/plixy/pkg/stats"
	"github.com/purini-to/plixy/pkg/trace"

	"github.com/pkg/errors"
	"github.com/purini-to/plixy/pkg/config"
	"github.com/purini-to/plixy/pkg/log"
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
		if err := stats.InitExporter(config.Global.Stats); err != nil {
			return errors.Wrap(err, "failed to create stats exporter")
		}
	}
	if config.Global.Trace.Enable {
		if err := trace.InitExporter(config.Global.Trace); err != nil {
			return errors.Wrap(err, "failed to create trace exporter")
		}
	}

	return nil
}
