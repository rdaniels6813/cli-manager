package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
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
		install, _ := cmd.Flags().GetBool("install")
		shellType := getShellType(cmd)

		switch shellType {
		case Zsh:
			handleZshCompletion(gen, install)
		case Bash:
			handleBashCompletion(gen, install)
		case Powershell:
			handlePowershellCompletion(gen, install, false)
		case PowershellCore:
			handlePowershellCompletion(gen, install, true)
		case Unknown:
			log.Fatal("Unknown shell, please specify your shell using flags")
		}
	},
}

const zshCompletionSnippet = "source <(cli-manager completion -g -z)"
const bashCompletionSnippet = "source <(cli-manager completion -g -b)"
const powershellCompletionSnippet = "Invoke-Expression $($(cli-manager.exe completion -g -p) -join \"`n\")"

func handleZshCompletion(generate bool, install bool) {
	if generate {
		var data []byte
		buf := bytes.NewBuffer(data)
		rootCmd.GenZshCompletion(buf)
		output, _ := ioutil.ReadAll(buf)
		fmt.Print(strings.TrimPrefix(string(output), "#"))
	} else if install {
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
			fmt.Printf("Wrote completion script to: %s", scriptPath)
		} else {
			fmt.Printf("Already installed in: %s", scriptPath)
		}
	} else {
		fmt.Printf("Add the following line to your .zshrc file:\n\n%s", zshCompletionSnippet)
	}
}

func handleBashCompletion(generate bool, install bool) {
	if generate {
		rootCmd.GenBashCompletion(os.Stdout)
	} else if install {
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
			fmt.Printf("Wrote completion script to: %s", scriptPath)
		} else {
			fmt.Printf("Already installed in: %s", scriptPath)
		}
	} else {
		fmt.Printf("Add the following line to your .bashrc or .profile file:\n\n%s", bashCompletionSnippet)
	}
}

func handlePowershellCompletion(generate bool, install bool, core bool) {
	if generate {
		rootCmd.GenPowerShellCompletion(os.Stdout)
	} else if install {
		scriptPath := getPowershellProfilePath(core)
		wrote, err := writeShellSnippet(powershellCompletionSnippet, scriptPath)
		if err != nil {
			log.Fatal(err)
		}
		if wrote {
			fmt.Printf("Wrote completion script to: %s", scriptPath)
		} else {
			fmt.Printf("Already installed in: %s", scriptPath)
		}
	} else {
		fmt.Printf("Add the following line to your $PROFILE file:\n\n%s", powershellCompletionSnippet)
	}
}

type shellType string

const (
	// Zsh enum for zsh shell
	Zsh shellType = "zsh"
	// Bash enum for the bash shell
	Bash shellType = "bash"
	// Powershell enum for powershell
	Powershell shellType = "powershell"
	// PowershellCore enum for powershell
	PowershellCore shellType = "pwsh"
	// Unknown enum for unknown or unsupported shell
	Unknown shellType = "unknown"
)

func getShellType(cmd *cobra.Command) shellType {
	zsh, _ := cmd.Flags().GetBool("zsh")
	powershell, _ := cmd.Flags().GetBool("powershell")
	bash, _ := cmd.Flags().GetBool("bash")
	powershellCore, _ := cmd.Flags().GetBool("pwsh")
	if zsh {
		return Zsh
	}
	if bash {
		return Bash
	}
	if powershell {
		return Powershell
	}
	if powershellCore {
		return PowershellCore
	}
	if os.Getenv("ZSH_NAME") != "" || os.Getenv("ZSH") != "" {
		return Zsh
	}
	if os.Getenv("BASH") != "" {
		return Bash
	}
	psModule := os.Getenv("PSModulePath")
	if psModule != "" {
		if strings.Contains(strings.ToLower(psModule), fmt.Sprintf("%spowershell%s", string(os.PathSeparator), string(os.PathSeparator))) {
			return PowershellCore
		}
		return Powershell
	}
	return Unknown
}

func getPowershellProfilePath(core bool) string {
	dir, _ := os.UserHomeDir()
	onedrivePath := os.Getenv("ONEDRIVE")
	myDocuments := ""
	if onedrivePath != "" {
		myDocuments = filepath.Join(onedrivePath, "Documents")
	}
	if _, err := os.Stat(myDocuments); os.IsNotExist(err) {
		myDocuments = filepath.Join(dir, "My Documents")
		if _, err = os.Stat(myDocuments); os.IsNotExist(err) {
			myDocuments = filepath.Join(dir, "Documents")
		}
	}
	if core {
		switch runtime.GOOS {
		case "windows":
			return filepath.Join(myDocuments, "PowerShell", "Microsoft.PowerShell_profile.ps1")
		case "linux":
			return filepath.Join(dir, ".config", "powershell", "Microsoft.PowerShell_profile.ps1")
		case "darwin":
			fmt.Print("This is untested on macos with powershell core")
			return filepath.Join(dir, ".config", "powershell", "Microsoft.PowerShell_profile.ps1")
		}

	}
	return filepath.Join(myDocuments, "PowerShell", "Microsoft.PowerShell_profile.ps1")
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
	rootCmd.AddCommand(completionCmd)
	completionCmd.Flags().BoolP("powershell", "p", false, "Generate powershell completion")
	completionCmd.Flags().BoolP("pwsh", "c", false, "Generate powershell core completion")
	completionCmd.Flags().BoolP("bash", "b", false, "Generate bash completion")
	completionCmd.Flags().BoolP("zsh", "z", false, "Generate zsh completion")
	completionCmd.Flags().BoolP("generate", "g", false, "Generate completion for shell specified by $SHELL and send to stdout")
	completionCmd.Flags().BoolP("install", "i", false, "Install the completion script into the users profile")
}
