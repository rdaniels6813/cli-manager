package shell

import (
	"fmt"
	"os"
	"strings"
)

type ShellType string

const (
	// Zsh enum for zsh shell
	Zsh ShellType = "zsh"
	// Bash enum for the bash shell
	Bash ShellType = "bash"
	// Powershell enum for powershell
	Powershell ShellType = "powershell"
	// PowershellCore enum for powershell
	PowershellCore ShellType = "pwsh"
	// Unknown enum for unknown or unsupported shell
	Unknown ShellType = "unknown"
)

func GetShellType(zsh, powershell, bash, powershellCore bool) ShellType {
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
