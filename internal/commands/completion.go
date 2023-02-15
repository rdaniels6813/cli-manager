package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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
			fmt.Println("Unknown shell, please specify your shell using flags")
			os.Exit(1)
		}
	},
}

const zshCompletionSnippet = "\nsource <(cli-manager completion -g -z)\n"
const bashCompletionSnippet = "\nsource <(cli-manager completion -g -b)\n"
const powershellCompletionSnippet = "\nInvoke-Expression $($(cli-manager.exe completion -g -p) -join \"`n\")\n"

func handleZshCompletion(generate bool, install bool) {
	switch {
	case generate:
		var data []byte
		buf := bytes.NewBuffer(data)
		err := rootCmd.GenZshCompletion(buf)
		if err != nil {
			fmt.Println(err)
		}
		output, _ := io.ReadAll(buf)
		fmt.Print(strings.TrimPrefix(string(output), "#"))
	case install:
		dir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		scriptPath := filepath.Join(dir, ".zshrc")
		wrote, err := writeShellSnippet(zshCompletionSnippet, scriptPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if wrote {
			fmt.Printf("Wrote completion script to: %s\n", scriptPath)
		} else {
			fmt.Printf("Completion already installed in: %s\n", scriptPath)
		}
	default:
		fmt.Printf("Add the following line to your .zshrc file:\n\n%s", zshCompletionSnippet)
	}
}

func handleBashCompletion(generate bool, install bool) {
	switch {
	case generate:
		err := rootCmd.GenBashCompletion(os.Stdout)
		if err != nil {
			fmt.Println(err)
		}
	case install:
		dir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		scriptPath := filepath.Join(dir, ".bashrc")
		wrote, err := writeShellSnippet(bashCompletionSnippet, scriptPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if wrote {
			fmt.Printf("Wrote completion script to: %s\n", scriptPath)
		} else {
			fmt.Printf("Completion already installed in: %s\n", scriptPath)
		}
	default:
		fmt.Printf("Add the following line to your .bashrc or .profile file:\n\n%s", bashCompletionSnippet)
	}
}

func handlePowershellCompletion(generate bool, install bool, core bool) {
	switch {
	case generate:
		err := rootCmd.GenPowerShellCompletion(os.Stdout)
		if err != nil {
			fmt.Println(err)
		}
	case install:
		scriptPath := getPowershellProfilePath(core)
		wrote, err := writeShellSnippet(powershellCompletionSnippet, scriptPath)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if wrote {
			fmt.Printf("Wrote completion script to: %s\n", scriptPath)
		} else {
			fmt.Printf("Completion already installed in: %s\n", scriptPath)
		}
	default:
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
		powershellPartial := fmt.Sprintf("%spowershell%s", string(os.PathSeparator), string(os.PathSeparator))
		if strings.Contains(strings.ToLower(psModule), powershellPartial) {
			return PowershellCore
		}
		return Powershell
	}
	return Unknown
}

func getPowershellProfilePath(core bool) string {
	dir, _ := os.UserHomeDir()
	myDocuments := ""
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
		err := os.MkdirAll(filepath.Dir(path), 0777)
		if err != nil {
			return false, err
		}
		f, err := os.Create(path)
		if err != nil {
			return false, err
		}
		defer f.Close()
		_, err = f.WriteString(snippet)
		return true, err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return false, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Contains(text, strings.TrimSpace(snippet)) {
			return false, nil
		}
	}
	_, err = f.WriteString(snippet)
	return true, err
}

func init() {
	rootCmd.AddCommand(completionCmd)
	completionCmd.Flags().BoolP("powershell", "p", false, "Generate powershell completion")
	completionCmd.Flags().BoolP("pwsh", "c", false, "Generate powershell core completion")
	completionCmd.Flags().BoolP("bash", "b", false, "Generate bash completion")
	completionCmd.Flags().BoolP("zsh", "z", false, "Generate zsh completion")
	completionCmd.Flags().BoolP("generate", "g", false,
		"Generate completion for shell specified by $SHELL and send to stdout")
	completionCmd.Flags().BoolP("install", "i", false, "Install the completion script into the users profile")
}
