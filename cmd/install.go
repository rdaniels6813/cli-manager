package cmd

import (
	"log"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"lab.bittrd.com/bittrd/cli-manager/nodeman"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a CLI",
	Long:  `Install a CLI application for local use.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lts := nodeman.GetLatestLTSNodeVersion()
		nodeManager := nodeman.NewManager(afero.NewOsFs())
		node := nodeManager.GetNode(lts)
		err := node.Npm("install", "-g", args[0])
		if err != nil {
			log.Fatalf("Failed to install npm CLI: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
