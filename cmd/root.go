package cmd

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/anton-fuji/dibo/internal/templates"
	"github.com/spf13/cobra"
)

var (
	version       = "dev"
	readBuildInfo = debug.ReadBuildInfo
)

func resolvedVersion() string {
	if version != "dev" {
		return version
	}

	buildInfo, ok := readBuildInfo()
	if !ok || buildInfo.Main.Path != "github.com/anton-fuji/dibo" {
		return version
	}

	buildVersion := buildInfo.Main.Version
	if buildVersion == "" || buildVersion == "(devel)" {
		return version
	}
	return buildVersion
}

var rootCmd = &cobra.Command{
	Use:     "dibo",
	Short:   "dibo is a CLI tool to generate .dockerignore files",
	Long:    `dibo (dockerignore boilerplates) helps you easily access .dockerignore boilerplates.`,
	Version: resolvedVersion(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// templateNames provides shell completion for template-name arguments,
// omitting names already present on the command line.
func templateNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	names, err := templates.List()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	used := make(map[string]struct{}, len(args))
	for _, a := range args {
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
