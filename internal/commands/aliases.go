package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rdaniels6813/cli-manager/internal/nodeman"
	"github.com/spf13/afero"

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
		case Unknown:
			fmt.Println("Unknown shell, please specify your shell using flags")
			os.Exit(1)
		}
	},
}

const zshAliasesSnippet = "\nsource <(cli-manager aliases -g -z)\n"
const bashAliasesSnippet = "\nsource <(cli-manager aliases -g -b)\n"
const powershellAliasesSnippet = "\nInvoke-Expression $($(cli-manager.exe aliases -g -p) -join \"`n\")\n"

func handleZshAliases(generate bool, install bool) {
	switch {
	case generate:
		manager := nodeman.NewManager(afero.NewOsFs())
		apps := manager.GetInstalledExecutables()
		for _, app := range apps {
			fmt.Printf("alias %s='cli-manager run %s'\n", app, app)
		}
	case install:
		dir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		scriptPath := filepath.Join(dir, ".zshrc")
		wrote, err := writeShellSnippet(zshAliasesSnippet, scriptPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if wrote {
			fmt.Printf("Wrote aliases script to: %s\n", scriptPath)
		} else {
			fmt.Printf("Aliases already installed in: %s\n", scriptPath)
		}
	default:
		fmt.Printf("Add the following line to your .zshrc file:\n\n%s", zshAliasesSnippet)
	}
}

func handleBashAliases(generate bool, install bool) {
	switch {
	case generate:
		manager := nodeman.NewManager(afero.NewOsFs())
		apps := manager.GetInstalledExecutables()
		for _, app := range apps {
			fmt.Printf("alias %s='cli-manager run %s'\n", app, app)
		}
	case install:
		dir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		scriptPath := filepath.Join(dir, ".bashrc")
		wrote, err := writeShellSnippet(bashAliasesSnippet, scriptPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if wrote {
			fmt.Printf("Wrote aliases script to: %s\n", scriptPath)
		} else {
			fmt.Printf("Aliases already installed in: %s\n", scriptPath)
		}
	default:
		fmt.Printf("Add the following line to your .bashrc or .profile file:\n\n%s", bashAliasesSnippet)
	}
}

func handlePowershellAliases(generate bool, install bool, core bool) {
	switch {
	case generate:
		manager := nodeman.NewManager(afero.NewOsFs())
		apps := manager.GetInstalledExecutables()
		for _, app := range apps {
			fmt.Printf("function %s { cli-manager.exe run %s @args }", app, app)
		}
	case install:
		scriptPath := getPowershellProfilePath(core)
		wrote, err := writeShellSnippet(powershellAliasesSnippet, scriptPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if wrote {
			fmt.Printf("Wrote aliases script to: %s\n", scriptPath)
		} else {
			fmt.Printf("Aliases already installed in: %s\n", scriptPath)
		}
	default:
		fmt.Printf("Add the following line to your $PROFILE file:\n\n%s", powershellAliasesSnippet)
	}
}

func init() {
	rootCmd.AddCommand(aliasesCmd)
	aliasesCmd.Flags().BoolP("generate", "g", false,
		"Generate completion for shell specified by $SHELL and send to stdout")
	aliasesCmd.Flags().BoolP("powershell", "p", false, "Generate powershell aliases")
	aliasesCmd.Flags().BoolP("pwsh", "c", false, "Generate powershell core aliases")
	aliasesCmd.Flags().BoolP("bash", "b", false, "Generate bash aliases")
	aliasesCmd.Flags().BoolP("zsh", "z", false, "Generate zsh aliases")
	aliasesCmd.Flags().BoolP("install", "i", false, "Install the aliases init to the default location.")
}
