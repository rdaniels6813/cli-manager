package cmd

import (
	"fmt"
	"os"

	"github.com/rdaniels6813/cli-manager/internal/nodeman"
	"github.com/spf13/afero"

	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall [appName]",
	Short: "Uninstall a CLI application",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appName := args[0]
		manager := nodeman.NewManager(afero.NewOsFs())
		app, err := manager.GetCLIApp(appName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		node := manager.GetNodeByPath(app.Path)
		err = node.Npm("remove", "-g", app.App)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = manager.MarkUninstalled(app.App)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
