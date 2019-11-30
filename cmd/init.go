package cmd

import (
	"github.com/pkg/errors"
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
