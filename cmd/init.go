package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/anton-fuji/dibo/internal/templates"
	"github.com/spf13/cobra"
)

const defaultOutput = ".dockerignore"

var (
	initForce       bool
	initAppend      bool
	initOutput      string
	initInteractive bool
	initRecursive   bool
)

var initCmd = &cobra.Command{
	Use:   "init [templates...]",
	Short: "Create a .dockerignore file",
	Long: `Create a .dockerignore file from one or more templates.

Use --interactive to choose templates from a numbered list instead of passing
their names as arguments. When no templates are provided, dibo detects the
project type and asks for confirmation before generating the file.`,
	Example: `  dibo init Go Secrets
  dibo init
  dibo init --recursive
  dibo init --interactive
  dibo init -i --output docker/.dockerignore`,
	Args: func(cmd *cobra.Command, args []string) error {
		if initInteractive {
			if len(args) > 0 {
				return fmt.Errorf("templates cannot be used with --interactive")
			}
			if initRecursive {
				return fmt.Errorf("--recursive cannot be used with --interactive")
			}
			return nil
		}
		if initRecursive && len(args) > 0 {
			return fmt.Errorf("--recursive cannot be used with explicit templates")
		}
		return nil
	},
	ValidArgsFunction: templateNames,
	RunE: func(cmd *cobra.Command, args []string) error {
		errOut := cmd.ErrOrStderr()

		if initForce && initAppend {
			return fmt.Errorf("--force and --append cannot be used together")
		}
		if initRecursive && (initInteractive || len(args) > 0) {
			return fmt.Errorf("--recursive requires automatic detection")
		}
		autoDetect := len(args) == 0 && !initInteractive
		if autoDetect {
			detections, err := detectProject(".", initRecursive)
			if err != nil {
				return err
			}
			if len(detections) == 0 {
				return fmt.Errorf("no supported project type detected in .")
			}
			args = printDetectionSummary(cmd.OutOrStdout(), detections)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Recommended templates: %s\n", strings.Join(args, ", "))
			confirmed, err := confirmGeneration(cmd.OutOrStdout(), cmd.InOrStdin(), initOutput)
			if err != nil {
				return err
			}
			if !confirmed {
				_, _ = fmt.Fprintln(cmd.OutOrStdout(), "Canceled.")
				return nil
			}
		}
		if initInteractive {
			var err error
			args, err = selectTemplates(cmd.OutOrStdout(), cmd.InOrStdin())
			if err != nil {
				return err
			}
		}

		body, missing, err := templates.Combine(args)
		if err != nil {
			return err
		}
		for _, m := range missing {
			_, _ = fmt.Fprintf(errOut, "Warning: template %q not found, skipping...\n", m)
		}

		if initAppend {
			if err := appendToFile(initOutput, body); err != nil {
				return err
			}
		} else {
			if !initForce {
				if _, statErr := os.Stat(initOutput); statErr == nil {
					return fmt.Errorf("%s already exists; use --force to overwrite or --append to add", initOutput)
				} else if !os.IsNotExist(statErr) {
					return fmt.Errorf("stat %s: %w", initOutput, statErr)
				}
			}
			if err := writeDockerignore(initOutput, body, true); err != nil {
				return err
			}
		}

		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s written successfully\n", initOutput)
		return nil
	},
}

func selectTemplates(out io.Writer, in io.Reader) ([]string, error) {
	names, err := templates.List()
	if err != nil {
		return nil, err
	}
	_, _ = fmt.Fprintln(out, "Available templates:")
	for i, name := range names {
		_, _ = fmt.Fprintf(out, "  %d. %s\n", i+1, name)
	}
	_, _ = fmt.Fprint(out, "Select templates by number or name (comma-separated): ")

	line, err := bufio.NewReader(in).ReadString('\n')
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("read template selection: %w", err)
	}
	if strings.TrimSpace(line) == "" {
		return nil, fmt.Errorf("no templates selected")
	}

	selected := make([]string, 0)
	seen := make(map[string]bool)
	for value := range strings.SplitSeq(line, ",") {
		value = strings.TrimSpace(value)
		if value == "" {
			return nil, fmt.Errorf("invalid empty template selection")
		}
		if index, parseErr := strconv.Atoi(value); parseErr == nil {
			if index < 1 || index > len(names) {
				return nil, fmt.Errorf("template number %d is out of range", index)
			}
			value = names[index-1]
		}
		_, canonical, readErr := templates.Read(value)
		if readErr != nil {
			return nil, fmt.Errorf("unknown template %q", value)
		}
		if !seen[canonical] {
			seen[canonical] = true
			selected = append(selected, canonical)
		}
	}
	return selected, nil
}

func confirmGeneration(out io.Writer, in io.Reader, output string) (bool, error) {
	_, _ = fmt.Fprintf(out, "Create %s? [y/N] ", output)
	line, err := bufio.NewReader(in).ReadString('\n')
	if err != nil && err != io.EOF {
		return false, fmt.Errorf("read confirmation: %w", err)
	}
	answer := strings.ToLower(strings.TrimSpace(line))
	return answer == "y" || answer == "yes", nil
}

func writeDockerignore(path, body string, force bool) error {
	if !force {
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("%s already exists; use --force to overwrite", path)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("stat %s: %w", path, err)
		}
	}
	content := "# Generated by dibo\n\n" + body
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// appendToFile appends body to path, surfacing both write and close errors.
func appendToFile(path, body string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	_, writeErr := f.WriteString(body)
	closeErr := f.Close()
	if writeErr != nil {
		return fmt.Errorf("append to %s: %w", path, writeErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close %s: %w", path, closeErr)
	}
	return nil
}

func init() {
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "overwrite the file if it already exists")
	initCmd.Flags().BoolVarP(&initAppend, "append", "a", false, "append to the file instead of overwriting")
	initCmd.Flags().StringVarP(&initOutput, "output", "o", defaultOutput, "output file path")
	initCmd.Flags().BoolVarP(&initInteractive, "interactive", "i", false, "select templates interactively")
	initCmd.Flags().BoolVarP(&initRecursive, "recursive", "r", false, "detect supported projects in nested directories")
	rootCmd.AddCommand(initCmd)
}
