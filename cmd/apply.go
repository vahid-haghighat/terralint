package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vahid-haghighat/terralint/cmd/internal"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Modifies the terraform files passed in",
	Long:  `Modifies the terraform files passed in'`,
	Args:  validateArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return internal.Apply(terraformPath)
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)
}
