package util

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

func GetCliManagerFolder(aos afero.Fs) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("There was an error determining the user home directory")
		os.Exit(1)
	}
	cliManagerDir := filepath.Join(homeDir, ".cli-manager")
	if _, err := aos.Stat(cliManagerDir); os.IsNotExist(err) {
		err = aos.MkdirAll(cliManagerDir, 0700)
		if err != nil {
			fmt.Println(err)
		}
	}
	return cliManagerDir
}
