package cmd

import (
	"fmt"

	"github.com/anton-fuji/dibo/internal/templates"
	"github.com/spf13/cobra"
)

var dumpCmd = &cobra.Command{
	Use:   "dump [templates...]",
	Short: "Dump templates to stdout",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		foundCnt := 0
		for _, arg := range args {
			filename := fmt.Sprintf("%s.dockerignore", arg)
			content, err := templates.FS.ReadFile(filename)
			if err != nil {
				fmt.Printf("Error: Template '%s' not found\n", arg)
				continue
			}
			fmt.Printf("### %s ###\n", arg)
			fmt.Println(string(content))
			foundCnt++
		}
		if foundCnt == 0 {
			return fmt.Errorf("no valid templates found")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
}
