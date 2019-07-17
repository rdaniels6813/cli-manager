package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate a completion script for the specified shell and output to stdout",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		powershell, _ := cmd.Flags().GetBool("powershell")
		if powershell {
			rootCmd.GenPowerShellCompletion(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
	completionCmd.Flags().BoolP("powershell", "p", false, "Generate powershell completion")
}
