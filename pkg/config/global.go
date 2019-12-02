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
}

func init() {
	viper.SetDefault("Port", "8080")
	viper.SetDefault("GraceTimeOut", "30s")
	viper.SetDefault("Debug", false)

	viper.BindEnv("Port", "PORT")
	viper.BindEnv("GraceTimeOut", "GRACE_TIME_OUT")
	viper.BindEnv("Debug", "DEBUG")
	viper.BindEnv("DatabaseDSN", "DATABASE_DSN")
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
