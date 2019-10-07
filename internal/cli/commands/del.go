package commands

import (
	"fmt"
	"net/url"
	"os"

	"github.com/craftyphotons/ten34/pkg/client"
	"github.com/spf13/cobra"
)

var delDatabase string

// NewDelCommand returns the "del" command
func NewDelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "del [options] <key>",
		Short: "Deletes the specified key",
		Run:   delCommandFunc,
	}

	cmd.Flags().StringVarP(&delDatabase, "database", "d", "", "URI of the database")

	return cmd
}

func delCommandFunc(cmd *cobra.Command, args []string) {
	key := args[0]

	db, err := url.Parse(delDatabase)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client, err := client.New(*db)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = client.Delete(key)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
