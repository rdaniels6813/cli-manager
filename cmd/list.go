package cmd

import (
	"fmt"
	"sort"

	"github.com/spf13/afero"
	"lab.bittrd.com/bittrd/cli-manager/nodeman"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of the installed commands",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		manager := nodeman.NewManager(afero.NewOsFs())
		apps := manager.GetInstalledExecutables()
		sort.Strings(apps)
		for _, app := range apps {
			fmt.Println(app)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
