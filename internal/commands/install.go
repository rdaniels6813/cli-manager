package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rdaniels6813/cli-manager/internal/nodeman"
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
		nodeVersion, err := nodeman.GetLatestNodeVersion(http.DefaultClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		nodeManager := nodeman.NewManager(afero.NewOsFs())
		node, err := nodeManager.GetNode(nodeVersion)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		output, _ := node.NpmView(args[0])
		engine, err := cmd.Flags().GetString("node-version")
		if err != nil || engine == "" {
			engine = output.Engines["node"]
		}
		version, err := nodeman.GetNodeVersionByRangeOrLTS(engine, http.DefaultClient)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		installNode, err := nodeManager.GetNode(version)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = installNode.Npm("install", "-g", args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = nodeManager.MarkInstalled(output.Name, output.GetBins(), installNode.BinPath(), args[0])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	installCmd.Flags().StringP("node-version", "n", "", "Specify a node version to use for install: --node-version 12.x")
	rootCmd.AddCommand(installCmd)
}
