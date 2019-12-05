package cmd

import (
	"context"
	"fmt"

	"github.com/purini-to/plixy/pkg/api"
	"github.com/purini-to/plixy/pkg/config"

	"github.com/spf13/viper"

	"github.com/purini-to/plixy/pkg/server"

	"github.com/purini-to/plixy/pkg/log"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

// StartOptions are the command flags
type StartOptions struct {
	port           uint
	configFilePath string
	watch          bool
}

// NewStartCmd creates a new http server command
func NewStartCmd(ctx context.Context) *cobra.Command {
	opts := &StartOptions{}

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts a plixy web server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunServerStart(ctx, opts)
		},
	}

	cmd.PersistentFlags().UintVarP(&opts.port, "port", "p", 8080, "The port on which to start the server")
	cmd.PersistentFlags().StringVarP(&opts.configFilePath, "config", "c", "", "Config file path")
	cmd.PersistentFlags().BoolVarP(&opts.watch, "watch", "", false, "Watch and reloading api definition files")
	cmd.PersistentFlags().BoolVarP(&opts.watch, "stats", "", false, "Enable stats exporter by prometheus")
	cmd.PersistentFlags().BoolVarP(&opts.watch, "trace", "", false, "Enable trace exporter by jaeger")

	viper.BindPFlag("Port", cmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("Watch", cmd.PersistentFlags().Lookup("watch"))
	viper.BindPFlag("Stats.Enable", cmd.PersistentFlags().Lookup("stats"))
	viper.BindPFlag("Trace.Enable", cmd.PersistentFlags().Lookup("trace"))

	return cmd
}

// RunServerStart is the run command to start plixy
func RunServerStart(ctx context.Context, ops *StartOptions) error {
	if err := initLog(); err != nil {
		return errors.Wrap(err, "failed initialize log")
	}
	if err := initConfig(ops.configFilePath); err != nil {
		return errors.Wrap(err, "failed initialize config")
	}
	if err := initExporter(); err != nil {
		return errors.Wrap(err, "failed initialize exporter")
	}

	if config.Global.Stats.Enable {
		exporter := server.NewStatsExporter()
		err := exporter.StartWithContext(ctx)
		if err != nil {
			return errors.Wrap(err, "failed start stats exporter server")
		}
		defer exporter.Close()
	}

	err := api.InitRepository(config.Global.DatabaseDSN)
	if err != nil {
		return errors.Wrap(err, "failed build repository")
	}
	defer api.Close()

	log.Info(fmt.Sprintf("Start plixy %s server...", version))

	s := server.New()

	ctx = ContextWithSignal(ctx)
	err = s.Start(ctx)
	if err != nil {
		return errors.Wrap(err, "could not start server")
	}
	defer s.Close()

	s.Wait()

	log.Info(fmt.Sprintf("Stop plixy %s server...", version))

	return nil
}
