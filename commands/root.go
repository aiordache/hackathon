package commands

import (
	dockercli "github.com/docker/cli/cli"
	"github.com/spf13/cobra"
)

// PluginName is the name of the plugin
const PluginName = "run"

// RootCommand returns the compose command with its child commands
func RootCommand() *cobra.Command { //nolint:gocyclo
	c := &cobra.Command{
		Short:            "Docker Run++",
		Use:              PluginName,
		TraverseChildren: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			if err := run(args); err != nil {
				return dockercli.StatusError{
					StatusCode: 1,
					Status:     err.Error(),
				}
			}
			return nil
		},
	}
	c.Flags().SetInterspersed(false)
	return c
}
