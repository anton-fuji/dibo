package cmd

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/anton-fuji/dibo/internal/templates"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		entries, err := fs.ReadDir(templates.FS, ".")
		if err != nil {
			return fmt.Errorf("failed to read templates: %w", err)
		}

		fmt.Println("Available templates:")
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".dockerignore") {
				name := strings.TrimSuffix(entry.Name(), ".dockerignore")
				fmt.Printf("- %s\n", name)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
