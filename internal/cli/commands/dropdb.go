package commands

import (
	"fmt"
	"net/url"
	"os"

	"github.com/craftyphotons/ten34/pkg/client"
	"github.com/spf13/cobra"
)

// NewDropdbCommand returns the "dropdb" command
func NewDropdbCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dropdb [options] <uri>",
		Short: "Deletes a database",
		Run:   dropdbCommandFunc,
	}

	return cmd
}

func dropdbCommandFunc(cmd *cobra.Command, args []string) {
	db, err := url.Parse(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client, err := client.New(*db)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = client.DropDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
