package commands

import (
	"context"

	"github.com/aiordache/hackathon/pkg/compose"
	"github.com/spf13/cobra"
)

type runOptions struct {
	image string
}

func RootCommand(ctx context.Context) *cobra.Command {
	opts := runOptions{}
	cmd := &cobra.Command{
		Use:    "run",
		Short:  "Docker run++",
		Hidden: true,
		Args:   cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.image = args[0]
			return runRun(ctx, opts)
		},
	}
	return cmd
}

func runRun(ctx context.Context, opts runOptions) error {
	image := opts.image
	if ok, err := compose.IsDockerComposeImage(ctx, image, ""); err == nil && ok {
		// do something
	}
	return nil
}
