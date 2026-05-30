package cmd

import (
	"fmt"
	"strings"

	"github.com/anton-fuji/dibo/internal/templates"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:               "search [keyword]",
	Short:             "Search available templates by keyword",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: templateNames,
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()
		keyword := strings.ToLower(args[0])

		names, err := templates.List()
		if err != nil {
			return err
		}

		matched := make([]string, 0, len(names))
		for _, n := range names {
			if strings.Contains(strings.ToLower(n), keyword) {
				matched = append(matched, n)
			}
		}

		if len(matched) == 0 {
			_, _ = fmt.Fprintf(out, "No templates matching %q.\n", args[0])
			return nil
		}
		for _, n := range matched {
			_, _ = fmt.Fprintln(out, n)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
