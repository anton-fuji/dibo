package cmd

import (
	"fmt"
	"io/fs"
	"sort"
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

		var tmplName []string
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".dockerignore") {
				name := strings.TrimSuffix(entry.Name(), ".dockerignore")
				tmplName = append(tmplName, name)
			}
		}

		if len(tmplName) == 0 {
			fmt.Println("No templates available.")
			return nil
		}

		sort.Strings(tmplName)
		fmt.Println("\nAvailable templates:")
		fmt.Println()
		for _, name := range tmplName {
			fmt.Printf("%s\n", name)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
