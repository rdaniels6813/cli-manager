package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

func main() {
	NewApp().getNodeBaseFolder()
}

// App the base application with boundaries
type App struct {
	os afero.Fs
}

// NewApp constructor for cli application type
func NewApp() *App {
	factory := new(App)
	factory.os = afero.NewOsFs()
	return factory
}

func (a *App) getNodeBaseFolder() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("There was an error determining the user home directory")
	}
	cliManagerDir := filepath.Join(homeDir, ".cli-manager", "node")
	if _, err := a.os.Stat(cliManagerDir); os.IsNotExist(err) {
		a.os.MkdirAll(cliManagerDir, 0700)
	}
	return cliManagerDir
}
