package cmd

import (
	"fmt"
	"os"

	"github.com/rdaniels6813/cli-manager/internal/version"
	"github.com/spf13/cobra"
)

// rootCmd represents the install command
var rootCmd = &cobra.Command{
	Use:     "cli-manager",
	Short:   "Manage installation & use of various CLI applications",
	Long:    ``,
	Version: version.GetModuleVersion(),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
}
