package commands

import (
	"fmt"
	"net/url"
	"os"

	"github.com/craftyphotons/ten34/pkg/client"
	"github.com/spf13/cobra"
)

var putDatabase string

// NewPutCommand returns the "put" command
func NewPutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "put [options] <key> <value>",
		Short: "Writes the specified key",
		Run:   putCommandFunc,
	}

	cmd.Flags().StringVarP(&putDatabase, "database", "d", "", "URI of the database")

	return cmd
}

func putCommandFunc(cmd *cobra.Command, args []string) {
	key := args[0]
	val := args[1]

	db, err := url.Parse(putDatabase)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client, err := client.New(*db)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = client.Put(key, val)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
