package commands

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/aiordache/hackathon/pkg/compose"
	"github.com/spf13/cobra"
)

type runOptions struct {
	imageName string
}

func RunCommand(ctx context.Context) *cobra.Command {
	opts := runOptions{}
	cmd := &cobra.Command{
		Use:   "run [OPTIONS]",
		Short: "Docker run++",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.imageName = args[0]
			return runRun(ctx, opts)
		},
	}
	return cmd
}

func runRun(ctx context.Context, opts runOptions) error {
	image := opts.imageName
	fmt.Printf("Running %s\n", image)
	b, err := compose.EmbeddedCompose(ctx, image)
	if err != nil {
		return fmt.Errorf("not a compose image: %s", err.Error())
	}
	cmd := exec.Command("docker", "compose", "-f", "-", "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = bytes.NewReader(b)
	return cmd.Start()
}
