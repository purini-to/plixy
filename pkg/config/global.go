package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var Global = &global{}

type global struct {
	Debug        bool
	Port         uint
	GraceTimeOut string
	DatabaseDSN  string
	Watch        bool
}

func init() {
	viper.SetDefault("Port", "8080")
	viper.SetDefault("GraceTimeOut", "30s")
	viper.SetDefault("Debug", false)
	viper.SetDefault("DatabaseDSN", "")
	viper.SetDefault("Watch", false)

	viper.BindEnv("Port", "PLIXY_PORT")
	viper.BindEnv("GraceTimeOut", "PLIXY_GRACE_TIME_OUT")
	viper.BindEnv("Debug", "PLIXY_DEBUG")
	viper.BindEnv("DatabaseDSN", "PLIXY_DATABASE_DSN")
	viper.BindEnv("Watch", "PLIXY_WATCH")
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
