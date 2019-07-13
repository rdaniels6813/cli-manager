package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mholt/archiver"
	"github.com/spf13/afero"
)

func main() {
	app := NewApp()
	binaryPath := app.downloadNodeBinary("10.16.0")
	log.Printf("Binary downloaded: %v", binaryPath)
	outputPath := app.unpackNodeBinary(binaryPath, "10.16.0")
	log.Printf("Binary unpacked to: %v", outputPath)
}

// App the base application with boundaries
type App struct {
	os         afero.Fs
	httpClient *http.Client
}

// NewApp constructor for cli application type
func NewApp() *App {
	factory := new(App)
	factory.os = afero.NewOsFs()
	factory.httpClient = &http.Client{}
	return factory
}

func (a *App) getNodeURL(version string, os string, arch string) string {
	extension := ".tar.xz"
	if os == "windows" {
		os = "win"
		extension = ".zip"
	}
	if os == "darwin" {
		extension = ".tar.gz"
	}
	if arch == "amd64" {
		arch = "x64"
	}
	return fmt.Sprintf("https://nodejs.org/dist/v%s/node-v%s-%s-%s%s", version, version, os, arch, extension)
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

func (a *App) downloadNodeBinary(version string) string {
	nodeBaseFolder := a.getNodeBaseFolder()
	url := a.getNodeURL(version, runtime.GOOS, runtime.GOARCH)
	// Create the file
	nodeBinaryPath := filepath.Join(nodeBaseFolder, filepath.Base(url))
	if _, err := a.os.Stat(nodeBinaryPath); os.IsNotExist(err) {
		out, err := os.Create(nodeBinaryPath)
		if err != nil {
			log.Fatalf("Failed to create destination file for node binary: %v", err)
		}
		defer out.Close()

		// Get the data
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("Failed download node binary: %v", err)
		}
		defer resp.Body.Close()

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Fatalf("Failed write node binary: %v", err)
		}
	}
	return nodeBinaryPath
}

func (a *App) unpackNodeBinary(path string, version string) string {
	nodeFolder := a.getNodeBaseFolder()
	outputFolder := filepath.Join(nodeFolder, version)
	archiver.Unarchive(path, outputFolder)
	dirs, err := ioutil.ReadDir(outputFolder)
	if err != nil {
		log.Fatalf("Failed to read dir: %v", err)
	}
	dirName := filepath.Join(outputFolder, dirs[0].Name())
	tmpPath := filepath.Join(nodeFolder, "tmp")
	err = os.Rename(dirName, tmpPath)
	if err != nil {
		log.Fatalf("Failed to rename dir: %v", err)
	}
	os.Remove(outputFolder)
	err = os.Rename(tmpPath, outputFolder)
	return outputFolder
}
