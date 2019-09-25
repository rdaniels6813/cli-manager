package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/afero"
)

type ProfileHelper struct {
	FS afero.Fs
}

func (p *ProfileHelper) WriteProfileSnippet(snippet, path string) (bool, error) {
	if _, err := p.FS.Stat(path); os.IsNotExist(err) {
		err := p.FS.MkdirAll(filepath.Dir(path), 777)
		if err != nil {
			return false, err
		}
		f, err := p.FS.Create(path)
		defer f.Close()
		if err != nil {
			return false, err
		}
		_, err = f.WriteString(snippet)
		return true, err
	}
	profileBytes, err := afero.ReadFile(p.FS, path)
	if err != nil {
		return false, err
	}
	text := string(profileBytes)
	if strings.Contains(text, snippet) {
		return false, nil
	}
	f, err := p.FS.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	defer f.Close()
	if err != nil {
		return false, err
	}
	_, err = f.WriteString(fmt.Sprintf("%s\n", snippet))
	return true, err
}

func (p *ProfileHelper) GetPowershellProfilePath(core bool) string {
	dir, _ := os.UserHomeDir()
	myDocuments := ""
	if _, err := p.FS.Stat(myDocuments); os.IsNotExist(err) {
		myDocuments = filepath.Join(dir, "My Documents")
		if _, err = p.FS.Stat(myDocuments); os.IsNotExist(err) {
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
