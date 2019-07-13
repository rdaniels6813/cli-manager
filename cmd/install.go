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
	Run: func(cmd *cobra.Command, args []string) {
		nodeManager := nodeman.NewManager(afero.NewOsFs())
		node := nodeManager.GetNode("10.16.0")
		err := node.Node("-v")
		log.Println(err)
		err = node.Npm("-v")
		log.Println(err)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
