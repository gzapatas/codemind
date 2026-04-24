package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/gzapatas/codemind/internal/ollama"
	"github.com/spf13/cobra"
)

func newOllamaCmd() *cobra.Command {
	var baseURL string
	cmd := &cobra.Command{
		Use:     "oc",
		Aliases: []string{"ollama"},
		Short:   "Check Ollama connectivity and list available models",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			client := ollama.New(baseURL)
			if err := client.Ping(ctx); err != nil {
				return fmt.Errorf("Ollama unreachable: %w", err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Ollama reachable")

			models, err := client.ListModels(ctx)
			if err != nil {
				return fmt.Errorf("failed to list models: %w", err)
			}
			if len(models) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No models found")
				return nil
			}
			for _, m := range models {
				fmt.Fprintln(cmd.OutOrStdout(), m)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&baseURL, "base-url", "http://localhost:11434", "Ollama base URL")
	return cmd
}
