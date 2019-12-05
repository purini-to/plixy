package config

import (
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var Global = &global{}

type global struct {
	Debug         bool
	Port          uint
	GraceTimeOut  time.Duration
	DatabaseDSN   string
	Watch         bool
	WatchInterval time.Duration
	Stats         stats
	Trace         trace
}

func (g *global) IsObservable() bool {
	return g.Stats.Enable || g.Trace.Enable
}

type stats struct {
	Enable      bool
	Port        uint
	ServiceName string
}

type trace struct {
	Enable            bool
	AgentEndpoint     string
	CollectorEndpoint string
	ServiceName       string
}

func init() {
	viper.SetDefault("Port", "8080")
	viper.SetDefault("GraceTimeOut", 30*time.Second)
	viper.SetDefault("Debug", false)
	viper.SetDefault("DatabaseDSN", "")
	viper.SetDefault("Watch", false)
	viper.SetDefault("WatchInterval", 2*time.Second)
	viper.SetDefault("Stats.Enable", false)
	viper.SetDefault("Stats.Port", 9090)
	viper.SetDefault("Stats.ServiceName", "plixy")
	viper.SetDefault("Trace.Enable", false)
	viper.SetDefault("Trace.AgentEndpoint", "localhost:6831")
	viper.SetDefault("Trace.CollectorEndpoint", "http://localhost:14268/api/traces")
	viper.SetDefault("Trace.ServiceName", "plixy")

	viper.BindEnv("Port", "PLIXY_PORT")
	viper.BindEnv("GraceTimeOut", "PLIXY_GRACE_TIME_OUT")
	viper.BindEnv("Debug", "PLIXY_DEBUG")
	viper.BindEnv("DatabaseDSN", "PLIXY_DATABASE_DSN")
	viper.BindEnv("Watch", "PLIXY_WATCH")
	viper.BindEnv("WatchInterval", "PLIXY_WATCH_INTERVAL")
	viper.BindEnv("Stats.Enable", "PLIXY_STATS_ENABLE")
	viper.BindEnv("Stats.Port", "PLIXY_STATS_PORT")
	viper.BindEnv("Stats.ServiceName", "PLIXY_STATS_SERVICE_NAME")
	viper.BindEnv("Trace.Enable", "PLIXY_TRACE_ENABLE")
	viper.BindEnv("Trace.AgentEndpoint", "PLIXY_TRACE_AGENT_ENDPOINT")
	viper.BindEnv("Trace.CollectorEndpoint", "PLIXY_TRACE_COLLECTOR_ENDPOINT")
	viper.BindEnv("Trace.ServiceName", "PLIXY_TRACE_SERVICE_NAME")
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
