package commands

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"

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
	fmt.Printf("Running image %s\n", image)
	b, err := compose.EmbeddedCompose(ctx, image)
	if err != nil {
		if errors.Is(err, compose.ErrLayerNotFound) {
			// do regular run command
		}
		return err
	}
	cmd := exec.Command("docker", "compose", "-f", "-", "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = bytes.NewReader(b)
	return cmd.Start()
}
