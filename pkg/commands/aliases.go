package cmd

import (
	"fmt"
	"log"

	"github.com/rdaniels6813/cli-manager/pkg/aliases"
	"github.com/rdaniels6813/cli-manager/pkg/shell"

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
		zsh, _ := cmd.Flags().GetBool("zsh")
		powershell, _ := cmd.Flags().GetBool("powershell")
		bash, _ := cmd.Flags().GetBool("bash")
		powershellCore, _ := cmd.Flags().GetBool("pwsh")
		shellType := shell.GetShellType(zsh, powershell, bash, powershellCore)

		var err error
		switch shellType {
		case shell.Zsh:
			err = handleZshAliases(gen, install)
		case shell.Bash:
			err = handleBashAliases(gen, install)
		case shell.Powershell:
			err = handlePowershellAliases(gen, install, false)
		case shell.PowershellCore:
			err = handlePowershellAliases(gen, install, true)
		case shell.Unknown:
			log.Fatal("Unknown shell, please specify your shell using flags")
		}
		if err != nil {
			log.Fatalf("Error occurred handling aliases: %s", err)
		}
	},
}

func handleZshAliases(generate bool, install bool) error {
	generator := aliases.NewZshGenerator()
	switch {
	case generate:
		fmt.Println(generator.Generate())
		return nil
	case install:
		return generator.Install()
	default:
		fmt.Printf("Add the following line to your .zshrc file:\n\n%s", aliases.ZshAliasesSnippet)
		return nil
	}
}

func handleBashAliases(generate bool, install bool) error {
	generator := aliases.NewBashGenerator()
	switch {
	case generate:
		fmt.Println(generator.Generate())
		return nil
	case install:
		return generator.Install()
	default:
		fmt.Printf("Add the following line to your .bashrc or .profile file:\n\n%s", aliases.BashAliasesSnippet)
		return nil
	}
}

func handlePowershellAliases(generate bool, install bool, core bool) error {
	generator := aliases.NewPowershellGenerator(core)
	switch {
	case generate:
		fmt.Println(generator.Generate())
		return nil
	case install:
		return generator.Install()
	default:
		fmt.Printf("Add the following line to your $PROFILE file:\n\n%s", aliases.PowershellAliasesSnippet)
		return nil
	}
}

func init() {
	RootCmd.AddCommand(aliasesCmd)
	aliasesCmd.Flags().BoolP("generate", "g", false,
		"Generate completion for shell specified by $SHELL and send to stdout")
	aliasesCmd.Flags().BoolP("powershell", "p", false, "Generate powershell aliases")
	aliasesCmd.Flags().BoolP("pwsh", "c", false, "Generate powershell core aliases")
	aliasesCmd.Flags().BoolP("bash", "b", false, "Generate bash aliases")
	aliasesCmd.Flags().BoolP("zsh", "z", false, "Generate zsh aliases")
	aliasesCmd.Flags().BoolP("install", "i", false, "Install the aliases init to the default location.")
}
