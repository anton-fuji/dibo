package cmd

import (
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/anton-fuji/dibo/internal/project"
	"github.com/anton-fuji/dibo/internal/templates"
	"github.com/spf13/cobra"
)

var (
	detectWrite     bool
	detectForce     bool
	detectOutput    string
	detectRecursive bool
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
		detections, err := detectProject(dir, detectRecursive)
		if err != nil {
			return err
		}
		if len(detections) == 0 {
			return fmt.Errorf("no supported project type detected in %s", dir)
		}

		recommended := printDetectionSummary(cmd.OutOrStdout(), detections)
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

func detectProject(dir string, recursive bool) ([]project.Detection, error) {
	if recursive {
		return project.DetectWithEvidenceRecursive(dir)
	}
	return project.DetectWithEvidence(dir)
}

func recommendedTemplates(detections []project.Detection) []string {
	recommended := []string{"Common"}
	seen := map[string]bool{"Common": true}
	for _, detection := range detections {
		if !seen[detection.Template] {
			recommended = append(recommended, detection.Template)
			seen[detection.Template] = true
		}
	}
	if !seen["Secrets"] {
		recommended = append(recommended, "Secrets")
	}
	return recommended
}

func printDetectionSummary(out io.Writer, detections []project.Detection) []string {
	found := make([]string, 0, len(detections))
	seen := make(map[string]bool, len(detections))
	for _, detection := range detections {
		if !seen[detection.Template] {
			found = append(found, detection.Template)
			seen[detection.Template] = true
		}
	}
	_, _ = fmt.Fprintf(out, "Detected: %s\n", strings.Join(found, ", "))
	_, _ = fmt.Fprintln(out, "Detection evidence:")
	for _, detection := range detections {
		label := detection.Template
		if detection.Directory != "" && detection.Directory != "." {
			label += " (" + detection.Directory + ")"
		}
		_, _ = fmt.Fprintf(out, "  %s: %s\n", label, strings.Join(detection.Signals, ", "))
	}
	return recommendedTemplates(detections)
}

func init() {
	detectCmd.Flags().BoolVarP(&detectWrite, "write", "w", false, "write the recommended .dockerignore")
	detectCmd.Flags().BoolVarP(&detectForce, "force", "f", false, "overwrite the output file when used with --write")
	detectCmd.Flags().StringVarP(&detectOutput, "output", "o", defaultOutput, "output file path when used with --write")
	detectCmd.Flags().BoolVarP(&detectRecursive, "recursive", "r", false, "detect supported projects in nested directories")
	rootCmd.AddCommand(detectCmd)
}
