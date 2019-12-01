package config

import (
	"github.com/purini-to/plixy/pkg/log"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Option func(*global)

func WithEnvPrefix(prefix string) Option {
	return func(c *global) {
		viper.SetEnvPrefix(prefix)
	}
}

func WithLoadFile(filePath string) Option {
	return func(c *global) {
		viper.SetConfigFile(filePath)
		if err := viper.ReadInConfig(); err != nil {
			log.Warn("Not found config file. Continue load environment variables.", zap.String("filePath", filePath))
		}
	}
}
