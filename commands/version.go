package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

const Version = "0.1.0"

type versionOptions struct {
	format string
	short  bool
}

func VersionCommand() *cobra.Command {
	opts := versionOptions{}
	cmd := &cobra.Command{
		Use:   "version [OPTIONS]",
		Short: "Show the version information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			runVersion(opts)
			return nil
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// overwrite parent PersistentPreRunE to avoid trying to load
			// compose file on version command if COMPOSE_FILE is set
			return nil
		},
	}
	// define flags for backward compatibility with com.docker.cli
	flags := cmd.Flags()
	flags.StringVarP(&opts.format, "format", "f", "", "Format the output. Values: [pretty | json]. (Default: pretty)")
	flags.BoolVar(&opts.short, "short", false, "Shows only Compose's version number")

	return cmd
}

func runVersion(opts versionOptions) {
	if opts.short {
		fmt.Println(strings.TrimPrefix(Version, "v"))
		return
	}
	if opts.format == "json" {
		fmt.Printf("{\"version\":%q}\n", Version)
		return
	}
	fmt.Printf("Docker Run version: %s\n", Version)
}
