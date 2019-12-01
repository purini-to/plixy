package cmd

import (
	"context"
	"fmt"

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

	viper.BindPFlag("Port", cmd.PersistentFlags().Lookup("port"))

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

	log.Info(fmt.Sprintf("Start plixy %s server...", version))

	s := server.New()

	ctx = ContextWithSignal(ctx)
	err := s.Start(ctx)
	if err != nil {
		return errors.Wrap(err, "could not start server")
	}
	defer s.Close()

	s.Wait()

	log.Info(fmt.Sprintf("Stop plixy %s server...", version))

	return nil
}
