package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Give instructions for enabling shell completion",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gen, _ := cmd.Flags().GetBool("generate")
		shellType := getShellType(cmd)

		switch shellType {
		case Zsh:
			handleZsh(gen)
		case Bash:
			handleBash(gen)
		case Powershell:
			handlePowershell(gen)
		}
	},
}

func handleZsh(generate bool) {
	if generate {
		var data []byte
		buf := bytes.NewBuffer(data)
		rootCmd.GenZshCompletion(buf)
		output, _ := ioutil.ReadAll(buf)
		fmt.Print(strings.TrimPrefix(string(output), "#"))
	} else {
		fmt.Printf("Add the following line to your .zshrc file:\n\nsource <(cli-manager completion -g -z)")
	}
}

func handleBash(generate bool) {
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

func handlePowershell(generate bool) {
	if generate {
		rootCmd.GenPowerShellCompletion(os.Stdout)
	} else {
		fmt.Printf("Add the following line to your $PROFILE file:\n\n")
	}
}

type shellType string

const (
	// Zsh enum for zsh shell
	Zsh shellType = "zsh"
	// Bash enum for the bash shell
	Bash shellType = "bash"
	// Powershell enum for powershell
	Powershell shellType = "pwsh"
	// Unknown enum for unknown or unsupported shell
	Unknown shellType = "unknown"
)

func getShellType(cmd *cobra.Command) shellType {
	zsh, _ := cmd.Flags().GetBool("zsh")
	powershell, _ := cmd.Flags().GetBool("powershell")
	bash, _ := cmd.Flags().GetBool("bash")
	if zsh || os.Getenv("ZSH_NAME") != "" {
		return Zsh
	}
	if bash || os.Getenv("BASH") != "" {
		return Bash
	}
	if powershell || os.Getenv("PSModulePath") != "" {
		return Powershell
	}
	return Unknown
}

func init() {
	rootCmd.AddCommand(completionCmd)
	completionCmd.Flags().BoolP("powershell", "p", false, "Generate powershell completion")
	completionCmd.Flags().BoolP("bash", "b", false, "Generate bash completion")
	completionCmd.Flags().BoolP("zsh", "z", false, "Generate zsh completion")
	completionCmd.Flags().BoolP("generate", "g", false, "Generate completion for shell specified by $SHELL and send to stdout")
}
