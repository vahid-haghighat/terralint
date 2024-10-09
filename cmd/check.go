package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vahid-haghighat/terralint/cmd/internal"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the terraform file/files for linter rules",
	Long: `Checks the terraform file/files for the linter rules and returns a list of
locations where any of the builtin rules are violated.`,
	Args: validateArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return internal.Check(terraformPath)
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
