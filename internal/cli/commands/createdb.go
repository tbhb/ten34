package commands

import (
	"fmt"
	"net/url"
	"os"

	"github.com/craftyphotons/ten34/pkg/client"
	"github.com/spf13/cobra"
)

// NewCreatedbCommand returns the "createdb" command
func NewCreatedbCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "createdb [options] <uri>",
		Short: "Creates a database",
		Run:   createdbCommandFunc,
	}

	return cmd
}

func createdbCommandFunc(cmd *cobra.Command, args []string) {
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

	err = client.CreateDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
