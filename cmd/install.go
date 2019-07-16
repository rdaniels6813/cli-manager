package cmd

import (
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
		nodeVersion := nodeman.GetLatestNodeVersion()
		nodeManager := nodeman.NewManager(afero.NewOsFs())
		node := nodeManager.GetNode(nodeVersion)
		output, _ := node.NpmView(args[0])
		engine := output.Engines["node"]
		version := nodeman.GetNodeVersionByRangeOrLTS(engine)
		installNode := nodeManager.GetNode(version)
		installNode.Npm("install", "-g", "eslint")
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
