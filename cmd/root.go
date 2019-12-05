package cmd

import (
	"context"

	"github.com/purini-to/plixy/pkg/config"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

var debug bool

// NewRootCmd creates a new instance of the root command
func NewRootCmd() *cobra.Command {
	ctx := context.Background()

	cmd := &cobra.Command{
		Use:   "plixy",
		Short: "Plixy is an API Gateway",
		Long: `
This is a lightweight API Gateway.`,
		Version: config.Version,
	}

	cmd.PersistentFlags().BoolVarP(&debug, "debug", "", false, "Launch in debug mode if there is a flag")

	viper.BindPFlag("Debug", cmd.PersistentFlags().Lookup("debug"))

	cmd.AddCommand(NewStartCmd(ctx))

	return cmd
}
