package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
		install, _ := cmd.Flags().GetBool("install")
		shellType := getShellType(cmd)

		switch shellType {
		case Zsh:
			handleZshAliases(gen, install)
		case Bash:
			handleBashAliases(gen, install)
		case Powershell:
			handlePowershellAliases(gen, install, false)
		case PowershellCore:
			handlePowershellAliases(gen, install, true)
		}
	},
}

const zshAliasesSnippet = "source <(cli-manager aliases -g -z)"
const bashAliasesSnippet = "source <(cli-manager aliases -g -b)"
const powershellAliasesSnippet = "Invoke-Expression $($(cli-manager.exe aliases -g -p) -join \"`n\")"

func handleZshAliases(generate bool, install bool) {
	if generate {
		manager := nodeman.NewManager(afero.NewOsFs())
		apps := manager.GetInstalledApps()
		for _, app := range apps {
			fmt.Printf("alias %s='cli-manager run %s'\n", app, app)
		}
	} else if install {
		dir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		scriptPath := filepath.Join(dir, ".zshrc")
		wrote, err := writeShellSnippet(zshAliasesSnippet, scriptPath)
		if err != nil {
			log.Fatal(err)
		}
		if wrote {
			fmt.Printf("Wrote aliases script to: %s", scriptPath)
		} else {
			fmt.Printf("Already installed in: %s", scriptPath)
		}
	} else {
		fmt.Printf("Add the following line to your .zshrc file:\n\n%s", zshAliasesSnippet)
	}
}

func handleBashAliases(generate bool, install bool) {
	if generate {
		manager := nodeman.NewManager(afero.NewOsFs())
		apps := manager.GetInstalledApps()
		for _, app := range apps {
			fmt.Printf("alias %s='cli-manager run %s'\n", app, app)
		}
	} else if install {
		dir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		scriptPath := filepath.Join(dir, ".bashrc")
		wrote, err := writeShellSnippet(bashAliasesSnippet, scriptPath)
		if err != nil {
			log.Fatal(err)
		}
		if wrote {
			fmt.Printf("Wrote aliases script to: %s", scriptPath)
		} else {
			fmt.Printf("Already installed in: %s", scriptPath)
		}
	} else {
		fmt.Printf("Add the following line to your .bashrc or .profile file:\n\n%s", bashAliasesSnippet)
	}
}

func handlePowershellAliases(generate bool, install bool, core bool) {
	if generate {
		manager := nodeman.NewManager(afero.NewOsFs())
		apps := manager.GetInstalledApps()
		for _, app := range apps {
			fmt.Printf("function %s { cli-manager.exe run %s @args }", app, app)
		}
	} else if install {
		scriptPath := getPowershellProfilePath(core)
		wrote, err := writeShellSnippet(powershellAliasesSnippet, scriptPath)
		if err != nil {
			log.Fatal(err)
		}
		if wrote {
			fmt.Printf("Wrote aliases script to: %s", scriptPath)
		} else {
			fmt.Printf("Already installed in: %s", scriptPath)
		}
	} else {
		fmt.Printf("Add the following line to your $PROFILE file:\n\n%s", powershellAliasesSnippet)
	}
}

func init() {
	rootCmd.AddCommand(aliasesCmd)
	aliasesCmd.Flags().BoolP("generate", "g", false, "Generate completion for shell specified by $SHELL and send to stdout")
	aliasesCmd.Flags().BoolP("powershell", "p", false, "Generate powershell aliases")
	aliasesCmd.Flags().BoolP("pwsh", "c", false, "Generate powershell core aliases")
	aliasesCmd.Flags().BoolP("bash", "b", false, "Generate bash aliases")
	aliasesCmd.Flags().BoolP("zsh", "z", false, "Generate zsh aliases")
	aliasesCmd.Flags().BoolP("install", "i", false, "Install the aliases init to the default location.")
}
