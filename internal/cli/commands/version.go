package commands

import "github.com/spf13/cobra"

// NewVersionCommand returns the "version" command
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints version information",
		Run:   versionCommandFunc,
	}

	return cmd
}

func versionCommandFunc(cmd *cobra.Command, args []string) {
}
