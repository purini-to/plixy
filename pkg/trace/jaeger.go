package trace

import (
	"context"

	"github.com/pkg/errors"
	"github.com/purini-to/plixy/pkg/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"contrib.go.opencensus.io/exporter/jaeger"
)

type JaegerExporter struct {
	*jaeger.Exporter
}

func (j *JaegerExporter) Start(ctx context.Context) error {
	return nil
}

func (j *JaegerExporter) Close() error {
	return nil
}

type JaegerOption struct {
	ServiceName       string
	CollectorEndpoint string
	AgentEndpoint     string
}

func NewJaegerExporter(opt *JaegerOption) (*JaegerExporter, error) {
	exporter, err := jaeger.NewExporter(jaeger.Options{
		CollectorEndpoint: opt.CollectorEndpoint,
		AgentEndpoint:     opt.AgentEndpoint,
		OnError: func(err error) {
			log.GetLogger().WithOptions(zap.AddStacktrace(zapcore.PanicLevel)).Error("Error jaeger exporter", zap.Error(err))
		},
		Process: jaeger.Process{
			ServiceName: opt.ServiceName,
		},
		BufferMaxCount: 10,
	})

	if err != nil {
		return nil, errors.Wrap(err, "failed to create jaeger exporter")
	}

	return &JaegerExporter{
		Exporter: exporter,
	}, nil
}
