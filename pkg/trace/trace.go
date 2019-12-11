package trace

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/pkg/errors"
	"github.com/purini-to/plixy/pkg/config"
	"github.com/purini-to/plixy/pkg/log"
	"go.opencensus.io/trace"
)

var exporter Exporter

type Exporter interface {
	trace.Exporter
	Start(context.Context) error
	Close() error
}

func InitExporter(conf config.Trace) error {
	if !conf.Enable {
		return nil
	}

	switch conf.Name {
	case "jaeger":
		log.Debug("Jaeger trace exporter chosen")
		exp, err := NewJaegerExporter(&JaegerOption{
			ServiceName:       conf.ServiceName,
			CollectorEndpoint: conf.CollectorEndpoint,
			AgentEndpoint:     conf.AgentEndpoint,
		})
		if err != nil {
			return errors.Wrap(err, "failed initialize jaeger exporter")
		}
		exporter = exp
	default:
		return errors.New(fmt.Sprintf("The selected name is not supported to trace exporter. name: %s", conf.Name))
	}

	trace.RegisterExporter(exporter)
	if conf.SamplingFraction <= 0 {
		log.Debug("The selected trace sampler by always")
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	} else {
		log.Debug("The selected trace sampler by probability", zap.Float64("fraction", 1/conf.SamplingFraction))
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(1 / conf.SamplingFraction)})
	}

	return nil
}

func Start(ctx context.Context) error {
	if exporter != nil {
		return exporter.Start(ctx)
	}
	return nil
}

func Close() error {
	if exporter != nil {
		return exporter.Close()
	}
	return nil
}
