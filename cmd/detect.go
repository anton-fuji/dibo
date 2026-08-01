package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/anton-fuji/dibo/internal/project"
	"github.com/anton-fuji/dibo/internal/templates"
	"github.com/spf13/cobra"
)

var (
	detectWrite  bool
	detectForce  bool
	detectOutput string
)

var detectCmd = &cobra.Command{
	Use:   "detect [directory]",
	Short: "Detect project types and recommend templates",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) == 1 {
			dir = args[0]
		}
		found, err := project.Detect(dir)
		if err != nil {
			return err
		}
		if len(found) == 0 {
			return fmt.Errorf("no supported project type detected in %s", dir)
		}

		recommended := append([]string{"Common"}, found...)
		recommended = append(recommended, "Secrets")
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Detected: %s\n", strings.Join(found, ", "))
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Recommended templates: %s\n", strings.Join(recommended, ", "))

		if !detectWrite {
			_, _ = fmt.Fprintln(cmd.OutOrStdout())
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Create it with: %s\n", initCommand(dir, recommended))
			return nil
		}
		output := detectOutput
		if !filepath.IsAbs(output) {
			output = filepath.Join(dir, output)
		}
		body, _, err := templates.Combine(recommended)
		if err != nil {
			return err
		}
		if err := writeDockerignore(output, body, detectForce); err != nil {
			return err
		}
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s written successfully\n", output)
		return nil
	},
}

func initCommand(dir string, templates []string) string {
	command := "dibo init " + strings.Join(templates, " ")
	if filepath.Clean(dir) != "." {
		output := filepath.Join(dir, defaultOutput)
		if strings.ContainsAny(output, " \t\n\"'\\") {
			output = strconv.Quote(output)
		}
		command += " --output " + output
	}
	return command
}

func init() {
	detectCmd.Flags().BoolVarP(&detectWrite, "write", "w", false, "write the recommended .dockerignore")
	detectCmd.Flags().BoolVarP(&detectForce, "force", "f", false, "overwrite the output file when used with --write")
	detectCmd.Flags().StringVarP(&detectOutput, "output", "o", defaultOutput, "output file path when used with --write")
	rootCmd.AddCommand(detectCmd)
}
