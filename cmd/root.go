package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/anton-fuji/dibo/internal/templates"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
)

var rootCmd = &cobra.Command{
	Use:     "dibo",
	Short:   "dibo is a CLI tool to generate .dockerignore files",
	Long:    `dibo (dockerignore boilerplates) helps you easily access .dockerignore boilerplates.`,
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func templateNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	names, err := templates.List()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	used := make(map[string]struct{}, len(args))
	for _, a := range names {
		used[strings.ToLower(a)] = struct{}{}
	}
	out := make([]string, 0, len(names))
	for _, n := range names {
		if _, dup := used[strings.ToLower(n)]; dup {
			continue
		}
		out = append(out, n)
	}

	return out, cobra.ShellCompDirectiveNoFileComp
}

func init() {
	// setup flags
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "%s\n" .Version}}`)
}
