package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anton-fuji/dibo/internal/project"
	"github.com/spf13/cobra"
)

var checkFile string

var requiredPatterns = map[string][]string{
	"Go":     {"bin/", "*.test"},
	"Node":   {"node_modules/", ".env"},
	"Python": {"__pycache__/", ".venv/", ".env"},
	"Ruby":   {"vendor/bundle/", "log/*.log", ".env"},
	"Rust":   {"target/"},
	"Java":   {"target/", "build/"},
	"PHP":    {"vendor/", ".env"},
	"dotNet": {"bin/", "obj/"},
}

var secretPatterns = []string{".env", "*.pem", "*.key"}

var checkCmd = &cobra.Command{
	Use:          "check [directory]",
	Short:        "Check a .dockerignore for common omissions",
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) == 1 {
			dir = args[0]
		}
		path := checkFile
		if !filepath.IsAbs(path) {
			path = filepath.Join(dir, path)
		}
		patterns, err := readPatterns(path)
		if err != nil {
			return err
		}

		issues := make([]string, 0)
		for _, pattern := range secretPatterns {
			if !patterns[pattern] {
				issues = append(issues, fmt.Sprintf("missing secret exclusion %q", pattern))
			}
		}
		detected, err := project.Detect(dir)
		if err != nil {
			return err
		}
		for _, template := range detected {
			for _, pattern := range requiredPatterns[template] {
				if !patterns[pattern] {
					issues = append(issues, fmt.Sprintf("%s project: missing %q", template, pattern))
				}
			}
		}

		out := cmd.OutOrStdout()
		if len(issues) == 0 {
			_, _ = fmt.Fprintf(out, "%s looks good.\n", path)
			return nil
		}
		_, _ = fmt.Fprintf(out, "%s: found %d issue(s):\n", path, len(issues))
		for _, issue := range issues {
			_, _ = fmt.Fprintf(out, "- %s\n", issue)
		}
		_, _ = fmt.Fprintln(out, "Run `dibo detect --write` to generate a baseline file.")
		return fmt.Errorf(".dockerignore check failed")
	},
}

func readPatterns(path string) (patterns map[string]bool, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer func() {
		if closeErr := f.Close(); err == nil && closeErr != nil {
			patterns = nil
			err = fmt.Errorf("close %s: %w", path, closeErr)
		}
	}()

	patterns = make(map[string]bool)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "!") {
			patterns[line] = true
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return patterns, nil
}

func init() {
	checkCmd.Flags().StringVarP(&checkFile, "file", "f", defaultOutput, "dockerignore file to check")
	rootCmd.AddCommand(checkCmd)
}
