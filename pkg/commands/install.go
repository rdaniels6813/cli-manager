package cmd

import (
	"fmt"

	"github.com/rdaniels6813/cli-manager/pkg/nodeman"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
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
		err := installNode.Npm("install", "-g", args[0])
		if err != nil {
			fmt.Println(err)
		}
		err = nodeManager.MarkInstalled(output.Name, output.Bin, installNode.BinPath(), args[0])
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(installCmd)
}
