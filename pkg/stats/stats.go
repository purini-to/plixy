package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/purini-to/plixy/pkg/config"
	"github.com/purini-to/plixy/pkg/log"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// Tags
var (
	KeyPath, _    = tag.NewKey("path")
	KeyApiName, _ = tag.NewKey("api_name")
)

// Measures
var (
	ApiDefinitionVersion = stats.Int64(
		"http/proxy/api/definition_version",
		"Proxy api definition version",
		stats.UnitDimensionless)
	ConcurrentRequestCount = stats.Int64(
		"http/proxy/concurrent_request_count",
		"Current count of HTTP requests",
		stats.UnitDimensionless)
)

// AllViews aggregates the metrics
var AllViews = []*view.View{
	// server
	{
		Name:        "http/proxy/request_count",
		Description: "Count of HTTP requests started",
		Measure:     ochttp.ServerRequestCount,
		Aggregation: view.Count(),
	},
	{
		Name:        "http/proxy/concurrent_request_count",
		Description: "Current count of HTTP requests",
		Measure:     ConcurrentRequestCount,
		Aggregation: view.Sum(),
	},
	{
		Name:        "http/proxy/request_count_by_method",
		Description: "Server request count by HTTP method",
		TagKeys:     []tag.Key{ochttp.Method},
		Measure:     ochttp.ServerRequestCount,
		Aggregation: view.Count(),
	},
	{
		Name:        "http/proxy/request_count_by_path",
		Description: "Server request count by url path",
		TagKeys:     []tag.Key{KeyPath},
		Measure:     ochttp.ServerRequestCount,
		Aggregation: view.Count(),
	},
	{
		Name:        "http/proxy/request_bytes",
		Description: "Size distribution of HTTP request body",
		Measure:     ochttp.ServerRequestBytes,
		Aggregation: ochttp.DefaultSizeDistribution,
	},
	{
		Name:        "http/proxy/response_count_by_status_code",
		Description: "Server response count by status code",
		TagKeys:     []tag.Key{ochttp.StatusCode},
		Measure:     ochttp.ServerLatency,
		Aggregation: view.Count(),
	},
	{
		Name:        "http/proxy/response_bytes",
		Description: "Size distribution of HTTP response body",
		Measure:     ochttp.ServerResponseBytes,
		Aggregation: ochttp.DefaultSizeDistribution,
	},
	{
		Name:        "http/proxy/latency",
		Description: "Latency distribution of HTTP requests",
		Measure:     ochttp.ServerLatency,
		Aggregation: ochttp.DefaultLatencyDistribution,
	},
	{
		Name:        "http/proxy/api/definition_version",
		Description: "Proxy api definition versions",
		Measure:     ApiDefinitionVersion,
		Aggregation: view.LastValue(),
	},
	// client
	{
		Name:        "http/client/sent_bytes",
		Measure:     ochttp.ClientSentBytes,
		Aggregation: ochttp.DefaultSizeDistribution,
		Description: "Total bytes sent in request body (not including headers), by HTTP method and response status",
		TagKeys:     []tag.Key{ochttp.KeyClientMethod, ochttp.KeyClientStatus},
	},

	{
		Name:        "http/client/received_bytes",
		Measure:     ochttp.ClientReceivedBytes,
		Aggregation: ochttp.DefaultSizeDistribution,
		Description: "Total bytes received in response bodies (not including headers but including error responses with bodies), by HTTP method and response status",
		TagKeys:     []tag.Key{ochttp.KeyClientMethod, ochttp.KeyClientStatus},
	},

	{
		Name:        "http/client/roundtrip_latency",
		Measure:     ochttp.ClientRoundtripLatency,
		Aggregation: ochttp.DefaultLatencyDistribution,
		Description: "End-to-end latency, by HTTP method and response status",
		TagKeys:     []tag.Key{ochttp.KeyClientMethod, ochttp.KeyClientStatus},
	},

	{
		Name:        "http/client/completed_count",
		Measure:     ochttp.ClientRoundtripLatency,
		Aggregation: view.Count(),
		Description: "Count of completed requests, by HTTP method and response status",
		TagKeys:     []tag.Key{ochttp.KeyClientMethod, ochttp.KeyClientStatus},
	},
	{
		Name:        "http/client/completed_count_by_api_name",
		Measure:     ochttp.ClientRoundtripLatency,
		Aggregation: view.Count(),
		Description: "Count of completed requests, by HTTP method and response status and api name",
		TagKeys:     []tag.Key{ochttp.KeyClientMethod, ochttp.KeyClientStatus, KeyApiName},
	},
}

var exporter Exporter

type Exporter interface {
	view.Exporter
	Start(context.Context) error
	Close() error
}

func InitExporter(conf config.Stats) error {
	if !conf.Enable {
		return nil
	}

	switch conf.Name {
	case "prometheus":
		log.Debug("Prometheus stats exporter chosen")
		exp, err := NewPrometheusExporter(&PrometheusOption{
			Namespace: conf.ServiceName,
			Port:      conf.Port,
		})
		if err != nil {
			return errors.Wrap(err, "failed initialize prometheus exporter")
		}
		exporter = exp
	default:
		return errors.New(fmt.Sprintf("The selected name is not supported to stats exporter. name: %s", conf.Name))
	}

	view.RegisterExporter(exporter)
	view.SetReportingPeriod(5 * time.Second)

	if err := view.Register(AllViews...); err != nil {
		return errors.Wrap(err, "failed to register server views")
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
