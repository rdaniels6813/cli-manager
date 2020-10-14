package cmd

import (
	"log"

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
		nodeVersion, err := nodeman.GetLatestNodeVersion()
		if err != nil {
			log.Fatal(err)
		}
		nodeManager := nodeman.NewManager(afero.NewOsFs())
		node, err := nodeManager.GetNode(nodeVersion)
		if err != nil {
			log.Fatal(err)
		}
		output, _ := node.NpmView(args[0])
		engine, err := cmd.Flags().GetString("node-version")
		if err != nil {
			engine = output.Engines["node"]
		}
		version, err := nodeman.GetNodeVersionByRangeOrLTS(engine)
		if err != nil {
			log.Fatal(err)
		}
		installNode, err := nodeManager.GetNode(version)
		if err != nil {
			log.Fatal(err)
		}
		err = installNode.Npm("install", "-g", args[0])
		if err != nil {
			log.Fatal(err)
		}
		err = nodeManager.MarkInstalled(output.Name, output.GetBins(), installNode.BinPath(), args[0])
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	installCmd.Flags().StringP("node-version", "n", "", "Specify a node version to use for install: --node-version 12.x")
	rootCmd.AddCommand(installCmd)
}
