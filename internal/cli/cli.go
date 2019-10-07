package cli

import (
	"fmt"
	"os"

	"github.com/craftyphotons/ten34/internal/cli/commands"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ten34",
	Short: "DNS as a key-value database",
}

func init() {
	rootCmd.AddCommand(
		commands.NewCreatedbCommand(),
		commands.NewDropdbCommand(),
		commands.NewDelCommand(),
		commands.NewGetCommand(),
		commands.NewPutCommand(),
		commands.NewVersionCommand(),
	)
}

// Start starts the ten34 CLI
func Start() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
