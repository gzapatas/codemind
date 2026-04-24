package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/codemind/project/internal/ingestion"
	"github.com/spf13/cobra"
)

func newChunkCmd() *cobra.Command {
	var file string
	var lang string
	cmd := &cobra.Command{
		Use:   "chunk",
		Short: "Chunk a source file using AST chunker (tree-sitter)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if file == "" {
				return fmt.Errorf("--file is required")
			}
			b, err := os.ReadFile(filepath.Clean(file))
			if err != nil {
				return err
			}
			chunks, err := ingestion.ChunkFile(b, lang)
			if err != nil {
				return err
			}
			if len(chunks) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no chunks found")
				return nil
			}
			for i, c := range chunks {
				fmt.Fprintf(cmd.OutOrStdout(), "#%d %s %s (%d-%d)\n", i+1, c.Kind, c.Name, c.StartLine, c.EndLine)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&file, "file", "", "Path to source file to chunk")
	cmd.Flags().StringVar(&lang, "lang", "go", "Language of the source file (default: go)")
	return cmd
}
