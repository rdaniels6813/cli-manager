package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rdaniels6813/cli-manager/pkg/shell"
	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Give instructions for enabling shell completion",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		gen, _ := cmd.Flags().GetBool("generate")
		install, _ := cmd.Flags().GetBool("install")
		zsh, _ := cmd.Flags().GetBool("zsh")
		powershell, _ := cmd.Flags().GetBool("powershell")
		bash, _ := cmd.Flags().GetBool("bash")
		powershellCore, _ := cmd.Flags().GetBool("pwsh")
		shellType := shell.GetShellType(zsh, powershell, bash, powershellCore)

		switch shellType {
		case shell.Zsh:
			handleZshCompletion(gen, install)
		case shell.Bash:
			handleBashCompletion(gen, install)
		case shell.Powershell:
			handlePowershellCompletion(gen, install, false)
		case shell.PowershellCore:
			handlePowershellCompletion(gen, install, true)
		case shell.Unknown:
			log.Fatal("Unknown shell, please specify your shell using flags")
		}
	},
}

const zshCompletionSnippet = "source <(cli-manager completion -g -z)\n"
const bashCompletionSnippet = "source <(cli-manager completion -g -b)\n"
const powershellCompletionSnippet = "Invoke-Expression $($(cli-manager.exe completion -g -p) -join \"`n\")\n"

func handleZshCompletion(generate bool, install bool) {
	switch {
	case generate:
		var data []byte
		buf := bytes.NewBuffer(data)
		err := RootCmd.GenZshCompletion(buf)
		if err != nil {
			fmt.Println(err)
		}
		output, _ := ioutil.ReadAll(buf)
		fmt.Print(strings.TrimPrefix(string(output), "#"))
		break
	case install:
		dir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		scriptPath := filepath.Join(dir, ".zshrc")
		wrote, err := writeShellSnippet(zshCompletionSnippet, scriptPath)
		if err != nil {
			log.Fatal(err)
		}
		if wrote {
			fmt.Printf("Wrote completion script to: %s\n", scriptPath)
		} else {
			fmt.Printf("Completion already installed in: %s\n", scriptPath)
		}
		break
	default:
		fmt.Printf("Add the following line to your .zshrc file:\n\n%s", zshCompletionSnippet)
	}
}

func handleBashCompletion(generate bool, install bool) {
	switch {
	case generate:
		err := RootCmd.GenBashCompletion(os.Stdout)
		if err != nil {
			fmt.Println(err)
		}
		break
	case install:
		dir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		scriptPath := filepath.Join(dir, ".bashrc")
		wrote, err := writeShellSnippet(bashCompletionSnippet, scriptPath)
		if err != nil {
			log.Fatal(err)
		}
		if wrote {
			fmt.Printf("Wrote completion script to: %s\n", scriptPath)
		} else {
			fmt.Printf("Completion already installed in: %s\n", scriptPath)
		}
		break
	default:
		fmt.Printf("Add the following line to your .bashrc or .profile file:\n\n%s", bashCompletionSnippet)
	}
}

func handlePowershellCompletion(generate bool, install bool, core bool) {
	switch {
	case generate:
		err := RootCmd.GenPowerShellCompletion(os.Stdout)
		if err != nil {
			fmt.Println(err)
		}
		break
	case install:
		profileHelper := shell.ProfileHelper{}
		scriptPath, err := profileHelper.GetPowershellProfilePath(core)
		if err != nil {
			log.Fatal(err)
		}
		wrote, err := writeShellSnippet(powershellCompletionSnippet, scriptPath)
		if err != nil {
			log.Fatal(err)
		}
		if wrote {
			fmt.Printf("Wrote completion script to: %s\n", scriptPath)
		} else {
			fmt.Printf("Completion already installed in: %s\n", scriptPath)
		}
		break
	default:
		fmt.Printf("Add the following line to your $PROFILE file:\n\n%s", powershellCompletionSnippet)
	}
}

func writeShellSnippet(snippet string, path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(path), 777)
		if err != nil {
			return false, err
		}
		f, err := os.Create(path)
		defer f.Close()
		if err != nil {
			return false, err
		}
		_, err = f.WriteString(snippet)
		return true, err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	defer f.Close()
	if err != nil {
		return false, err
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, snippet) {
			return false, nil
		}
	}
	_, err = f.WriteString(fmt.Sprintf("%s\n", snippet))
	return true, err
}

func init() {
	RootCmd.AddCommand(completionCmd)
	completionCmd.Flags().BoolP("powershell", "p", false, "Generate powershell completion")
	completionCmd.Flags().BoolP("pwsh", "c", false, "Generate powershell core completion")
	completionCmd.Flags().BoolP("bash", "b", false, "Generate bash completion")
	completionCmd.Flags().BoolP("zsh", "z", false, "Generate zsh completion")
	completionCmd.Flags().BoolP("generate", "g", false,
		"Generate completion for shell specified by $SHELL and send to stdout")
	completionCmd.Flags().BoolP("install", "i", false, "Install the completion script into the users profile")
}
