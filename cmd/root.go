package cmd

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "codemind",
		Short: "CodeMind CLI",
	}
	root.AddCommand(newRunCmd())
	root.AddCommand(newOllamaCmd())
	root.AddCommand(newChunkCmd())
	return root
}
