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
}

func init() {
	viper.SetDefault("Port", "8080")
	viper.SetDefault("GraceTimeOut", 30*time.Second)
	viper.SetDefault("Debug", false)
	viper.SetDefault("DatabaseDSN", "")
	viper.SetDefault("Watch", false)
	viper.SetDefault("WatchInterval", 2*time.Second)

	viper.BindEnv("Port", "PLIXY_PORT")
	viper.BindEnv("GraceTimeOut", "PLIXY_GRACE_TIME_OUT")
	viper.BindEnv("Debug", "PLIXY_DEBUG")
	viper.BindEnv("DatabaseDSN", "PLIXY_DATABASE_DSN")
	viper.BindEnv("Watch", "PLIXY_WATCH")
	viper.BindEnv("WatchInterval", "PLIXY_WATCH_INTERVAL")
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
