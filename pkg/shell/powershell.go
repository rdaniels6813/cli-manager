package shell

import (
	"fmt"
	"os"
	"path/filepath"
)

func (p *ProfileHelper) GetPowershellProfilePath(core bool) (string, error) {
	dir, _ := os.UserHomeDir()

	if core {
		switch p.GOOS {
		case "windows":
			myDocuments, err := p.getWindowsDocumentsDirectory(dir)
			if err != nil {
				return "", err
			}
			return filepath.Join(myDocuments, "PowerShell", "Microsoft.PowerShell_profile.ps1"), nil
		case "linux":
			return filepath.Join(dir, ".config", "powershell", "Microsoft.PowerShell_profile.ps1"), nil
		case "darwin":
			fmt.Print("This is untested on macos with powershell core")
			return filepath.Join(dir, ".config", "powershell", "Microsoft.PowerShell_profile.ps1"), nil
		}

	}
	myDocuments, err := p.getWindowsDocumentsDirectory(dir)
	if err != nil {
		return "", err
	}
	return filepath.Join(myDocuments, "PowerShell", "Microsoft.PowerShell_profile.ps1"), nil
}

func (p *ProfileHelper) getWindowsDocumentsDirectory(dir string) (string, error) {
	myDocuments := filepath.Join(dir, "My Documents")
	if _, err := p.FS.Stat(myDocuments); os.IsNotExist(err) {
		myDocuments = filepath.Join(dir, "Documents")
		if _, err := p.FS.Stat(myDocuments); os.IsNotExist(err) {
			return "", fmt.Errorf("Failed to find an existing profile directory for powershell")
		}
	}
	return myDocuments, nil
}
