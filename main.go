package main

import (
	"aiordache/hackathon/commands"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

var version = "0.0.0"

func main() {
	plugin.Run(func(dockerCli command.Cli) *cobra.Command {
		rootCmd := commands.RootCommand()
		originalPreRun := rootCmd.PersistentPreRunE
		rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
			if err := plugin.PersistentPreRunE(cmd, args); err != nil {
				return err
			}
			if originalPreRun != nil {
				return originalPreRun(cmd, args)
			}
			return nil
		}
		return rootCmd
	}, manager.Metadata{
		SchemaVersion: "0.1.0",
		Vendor:        "Docker Inc.",
		Version:       version,
	})
}
