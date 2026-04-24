package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run the ingestion indexer",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("run: not implemented yet")
			return nil
		},
	}
}
