package commands

import (
	"context"

	"github.com/spf13/cobra"
)

func RootCommand(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "v2",
		Short:  "Docker CLI  build/run v2",
		Hidden: false,
	}
	cmd.AddCommand(BuildCommand(ctx))
	cmd.AddCommand(RunCommand(ctx))
	return cmd
}
