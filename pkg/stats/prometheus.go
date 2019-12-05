package stats

import (
	"contrib.go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var PrometheusExporter *prometheus.Exporter

// Tags
var (
	KeyPath, _ = tag.NewKey("path")
)

// Measures
var (
	ApiDefinitionVersion = stats.Int64(
		"proxy/api/definition/version",
		"Proxy api definition versions",
		stats.UnitDimensionless)
)

// AllViews aggregates the metrics
var AllViews = []*view.View{
	// request
	{
		Name:        "http/proxy/request_count",
		Description: "Count of HTTP requests started",
		Measure:     ochttp.ServerRequestCount,
		Aggregation: view.Count(),
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
		Name:        "proxy/api/definition/version",
		Description: "Proxy api definition versions",
		Measure:     ApiDefinitionVersion,
		Aggregation: view.LastValue(),
	},
}
