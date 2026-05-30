package cmd

import (
	"fmt"

	"github.com/anton-fuji/dibo/internal/templates"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available templates",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()

		names, err := templates.List()
		if err != nil {
			return err
		}
		if len(names) == 0 {
			_, _ = fmt.Fprintln(out, "No templates available.")
			return nil
		}

		_, _ = fmt.Fprintln(out, "Available templates:")
		_, _ = fmt.Fprintln(out)
		for _, name := range names {
			_, _ = fmt.Fprintln(out, name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
