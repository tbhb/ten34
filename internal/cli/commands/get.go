package commands

import (
	"fmt"
	"net/url"
	"os"

	"github.com/craftyphotons/ten34/pkg/client"
	"github.com/spf13/cobra"
)

var getDatabase string

// NewGetCommand returns the "get" command
func NewGetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [options] <key>",
		Short: "Gets the specified key",
		Run:   getCommandFunc,
	}

	cmd.Flags().StringVarP(&getDatabase, "database", "d", "", "URI of the database")

	return cmd
}

func getCommandFunc(cmd *cobra.Command, args []string) {
	key := args[0]

	db, err := url.Parse(getDatabase)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client, err := client.New(*db)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	values, err := client.Get(key)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(values)
}
