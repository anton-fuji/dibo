package cmd

import (
	"fmt"
	"os"

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

func init() {
	// setup flags
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "%s\n" .Version}}`)
}
