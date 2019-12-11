package config

import (
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var Global = &global{}

type global struct {
	Debug               bool
	Port                uint
	GraceTimeOut        time.Duration
	DatabaseDSN         string
	Watch               bool
	WatchInterval       time.Duration
	DialTimeout         time.Duration
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     time.Duration
	Stats               Stats
	Trace               Trace
}

func (g *global) IsObservable() bool {
	return g.Stats.Enable || g.Trace.Enable
}

type Stats struct {
	Enable      bool
	Name        string
	Port        uint
	ServiceName string
}

type Trace struct {
	Enable            bool
	Name              string
	AgentEndpoint     string
	CollectorEndpoint string
	ServiceName       string
	SamplingFraction  float64
}

func init() {
	viper.SetDefault("Port", "8080")
	viper.SetDefault("GraceTimeOut", 30*time.Second)
	viper.SetDefault("Debug", false)
	viper.SetDefault("Watch", false)
	viper.SetDefault("WatchInterval", 2*time.Second)
	viper.SetDefault("DialTimeout", 30*time.Second)
	viper.SetDefault("MaxIdleConns", 512)
	viper.SetDefault("MaxIdleConnsPerHost", 128)
	viper.SetDefault("IdleConnTimeout", 90*time.Second)
	viper.SetDefault("Stats.Enable", false)
	viper.SetDefault("Stats.Name", "prometheus")
	viper.SetDefault("Stats.Port", 9090)
	viper.SetDefault("Stats.ServiceName", "plixy")
	viper.SetDefault("Trace.Enable", false)
	viper.SetDefault("Trace.Name", "jaeger")
	viper.SetDefault("Trace.ServiceName", "plixy")

	viper.BindEnv("Port", "PLIXY_PORT")
	viper.BindEnv("GraceTimeOut", "PLIXY_GRACE_TIME_OUT")
	viper.BindEnv("Debug", "PLIXY_DEBUG")
	viper.BindEnv("DatabaseDSN", "PLIXY_DATABASE_DSN")
	viper.BindEnv("Watch", "PLIXY_WATCH")
	viper.BindEnv("WatchInterval", "PLIXY_WATCH_INTERVAL")
	viper.BindEnv("DialTimeout", "PLIXY_DIAL_TIMEOUT")
	viper.BindEnv("MaxIdleConns", "PLIXY_MAX_IDLE_CONNS")
	viper.BindEnv("MaxIdleConnsPerHost", "PLIXY_MAX_IDLE_CONNS_PER_HOST")
	viper.BindEnv("IdleConnTimeout", "PLIXY_IDLE_CONN_TIMEOUT")
	viper.BindEnv("Stats.Enable", "PLIXY_STATS_ENABLE")
	viper.BindEnv("Stats.Name", "PLIXY_STATS_NAME")
	viper.BindEnv("Stats.Port", "PLIXY_STATS_PORT")
	viper.BindEnv("Stats.ServiceName", "PLIXY_STATS_SERVICE_NAME")
	viper.BindEnv("Trace.Enable", "PLIXY_TRACE_ENABLE")
	viper.BindEnv("Trace.Name", "PLIXY_TRACE_NAME")
	viper.BindEnv("Trace.AgentEndpoint", "PLIXY_TRACE_AGENT_ENDPOINT")
	viper.BindEnv("Trace.CollectorEndpoint", "PLIXY_TRACE_COLLECTOR_ENDPOINT")
	viper.BindEnv("Trace.ServiceName", "PLIXY_TRACE_SERVICE_NAME")
	viper.BindEnv("Trace.SamplingFraction", "PLIXY_TRACE_SAMPLING_FRACTION")
}

func Load(ops ...Option) error {
	viper.AutomaticEnv()

	var config global
	for _, o := range ops {
		o(&config)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return errors.Wrap(err, "failed unmarshal for config")
	}

	Global = &config
	return nil
}
