package cmd

import (
	"context"

	"github.com/purini-to/plixy/pkg/log"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

// NewStartCmd creates a new http server command
func NewStartCmd(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Starts a plixy web server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunServerStart(ctx)
		},
	}

	return cmd
}

// RunServerStart is the run command to start plixy
func RunServerStart(ctx context.Context) error {
	if err := initLog(); err != nil {
		return errors.Wrap(err, "Failed initialize log")
	}

	log.Info("Start plixy server...")
	return nil
}
