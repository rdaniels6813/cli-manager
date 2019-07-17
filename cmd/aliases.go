package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/afero"
	"lab.bittrd.com/bittrd/cli-manager/nodeman"

	"github.com/spf13/cobra"
)

// aliasesCmd represents the completion command
var aliasesCmd = &cobra.Command{
	Use:   "aliases",
	Short: "Creates aliases to all the installed apps",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gen, _ := cmd.Flags().GetBool("generate")
		shellType := getShellType(cmd)

		switch shellType {
		case Zsh:
			handleZshAliases(gen)
		case Bash:
			handleBashAliases(gen)
		case Powershell:
			handlePowershellAliases(gen)
		}
	},
}

func handleZshAliases(generate bool) {
	if generate {
		manager := nodeman.NewManager(afero.NewOsFs())
		apps := manager.GetInstalledApps()
		for _, app := range apps {
			fmt.Printf("alias %s='cli-manager run %s'", app, app)
		}
	} else {
		fmt.Printf("Add the following line to your .zshrc file:\n\nsource <(cli-manager aliases -g -z)")
	}
}

func handleBashAliases(generate bool) {
	if generate {
		var data []byte
		buf := bytes.NewBuffer(data)
		rootCmd.GenBashCompletion(buf)
		output, _ := ioutil.ReadAll(buf)
		log.Println(string(output))
	} else {
		fmt.Printf("Add the following line to your .bashrc or .profile file:\n\n")
	}
}

func handlePowershellAliases(generate bool) {
	if generate {
		rootCmd.GenPowerShellCompletion(os.Stdout)
	} else {
		fmt.Printf("Add the following line to your $PROFILE file:\n\n")
	}
}

func init() {
	rootCmd.AddCommand(aliasesCmd)
	aliasesCmd.Flags().BoolP("generate", "g", false, "Generate completion for shell specified by $SHELL and send to stdout")
	aliasesCmd.Flags().BoolP("powershell", "p", false, "Generate powershell aliases")
	aliasesCmd.Flags().BoolP("bash", "b", false, "Generate bash aliases")
	aliasesCmd.Flags().BoolP("zsh", "z", false, "Generate zsh aliases")
}
