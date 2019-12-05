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
	KeyPath, _    = tag.NewKey("path")
	KeyApiName, _ = tag.NewKey("api_name")
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
	// server
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
		Name:        "http/proxy/api/definition/version",
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
