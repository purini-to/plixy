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
}

type stats struct {
	Enable      bool
	Port        uint
	ServiceName string
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

	viper.BindEnv("Port", "PLIXY_PORT")
	viper.BindEnv("GraceTimeOut", "PLIXY_GRACE_TIME_OUT")
	viper.BindEnv("Debug", "PLIXY_DEBUG")
	viper.BindEnv("DatabaseDSN", "PLIXY_DATABASE_DSN")
	viper.BindEnv("Watch", "PLIXY_WATCH")
	viper.BindEnv("WatchInterval", "PLIXY_WATCH_INTERVAL")
	viper.BindEnv("Stats.Enable", "PLIXY_STATS_ENABLE")
	viper.BindEnv("Stats.Port", "PLIXY_STATS_PORT")
	viper.BindEnv("Stats.ServiceName", "PLIXY_STATS_SERVICENAME")
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
