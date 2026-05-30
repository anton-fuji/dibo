package cmd

import (
	"fmt"

	"github.com/anton-fuji/dibo/internal/templates"
	"github.com/spf13/cobra"
)

var dumpCmd = &cobra.Command{
	Use:               "dump [templates...]",
	Short:             "Dump templates to stdout",
	Args:              cobra.MinimumNArgs(1),
	ValidArgsFunction: templateNames,
	RunE: func(cmd *cobra.Command, args []string) error {
		out := cmd.OutOrStdout()
		errOut := cmd.ErrOrStderr()

		found := 0
		for _, arg := range args {
			content, canonical, err := templates.Read(arg)
			if err != nil {
				_, _ = fmt.Fprintf(errOut, "Error: template %q not found\n", arg)
				continue
			}
			_, _ = fmt.Fprintf(out, "### %s ###\n", canonical)
			_, _ = fmt.Fprintln(out, string(content))
			found++
		}

		if found == 0 {
			return fmt.Errorf("no valid templates found")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
}
