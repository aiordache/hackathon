package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type buildOptions struct {
	imageName       string
	composeFilePath string
}

func BuildCommand(ctx context.Context) *cobra.Command {
	opts := buildOptions{}
	cmd := &cobra.Command{
		Use:   "build [OPTIONS] PATH",
		Short: "Docker build++",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.composeFilePath = args[0]
			return runBuild(ctx, opts)
		},
	}
	cmd.Flags().StringVarP(&opts.imageName, "tag", "t", ".", "Image name/tag to build")
	return cmd
}

func runBuild(ctx context.Context, opts buildOptions) error {
	composeFile := opts.composeFilePath
	if !strings.HasSuffix(opts.composeFilePath, ".yaml") || !strings.HasSuffix(opts.composeFilePath, ".yml") {
		composeFile = filepath.Join(opts.composeFilePath, "docker-compose.yaml")
		if _, err := os.Stat(composeFile); err != nil {
			composeFile = filepath.Join(opts.composeFilePath, "docker-compose.yml")
			if _, err := os.Stat(composeFile); err != nil {
				return fmt.Errorf("no compose file found in build context: %s", err.Error())
			}
		}
	}
	dockerfile := []byte(`FROM scratch
LABEL composehackv1=true
ADD ` + filepath.Base(composeFile) + ` .
`)
	args := []string{"build", "--push", "-t", opts.imageName, "--annotation", "composehackv1=true", filepath.Dir(composeFile)}
	os.WriteFile(filepath.Join(filepath.Dir(composeFile), "Dockerfile"), dockerfile, 0644)
	defer os.Remove(filepath.Join(filepath.Dir(composeFile), "Dockerfile"))
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
