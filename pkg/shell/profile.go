package shell

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func WriteProfileSnippet(snippet, path string) (bool, error) {
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

func GetPowershellProfilePath(core bool) string {
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
