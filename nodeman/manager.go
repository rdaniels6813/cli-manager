package nodeman

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/afero"

	"github.com/mholt/archiver"
)

// Manager struct for managing node binaries leveraged by the cli manager
// Use NewManager to create an instance of this object
type Manager struct {
	os   afero.Fs
	http *http.Client
}

// NewManager constructor for default manager with the specified node version
func NewManager(os afero.Fs) *Manager {
	return &Manager{os: os}
}

// GetNode Ensures the node binary is downloaded & extracted and returns the path to the bin folder
func (m *Manager) GetNode(version string) Node {
	destinationPath := m.getNodeOutputFolder(version)
	if _, err := m.os.Stat(destinationPath); os.IsNotExist(err) {
		archivePath := m.downloadNodeArchive(version)
		m.unpackNodeArchive(archivePath, version)
	}
	return newNode(destinationPath)
}

func (m *Manager) getNodeURL(version string, os string, arch string) string {
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

func (m *Manager) getNodeBaseFolder() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("There was an error determining the user home directory")
	}
	cliManagerDir := filepath.Join(homeDir, ".cli-manager", "node")
	if _, err := m.os.Stat(cliManagerDir); os.IsNotExist(err) {
		m.os.MkdirAll(cliManagerDir, 0700)
	}
	return cliManagerDir
}

func (m *Manager) downloadNodeArchive(version string) string {
	nodeBaseFolder := m.getNodeBaseFolder()
	url := m.getNodeURL(version, runtime.GOOS, runtime.GOARCH)
	// Create the file
	nodeBinaryPath := filepath.Join(nodeBaseFolder, filepath.Base(url))
	if _, err := m.os.Stat(nodeBinaryPath); os.IsNotExist(err) {
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
func (m *Manager) getNodeOutputFolder(version string) string {
	nodeFolder := m.getNodeBaseFolder()
	return filepath.Join(nodeFolder, version)
}

func (m *Manager) unpackNodeArchive(path string, version string) {
	outputFolder := m.getNodeOutputFolder(version)
	archiver.Unarchive(path, outputFolder)
	dirs, err := ioutil.ReadDir(outputFolder)
	if err != nil {
		log.Fatalf("Failed to read dir: %v", err)
	}
	dirName := filepath.Join(outputFolder, dirs[0].Name())
	tmpPath, err := ioutil.TempDir(m.getNodeBaseFolder(), version)
	os.Remove(tmpPath)
	if err != nil {
		log.Fatalf("Failed to get temp dir: %v", err)
	}
	err = os.Rename(dirName, tmpPath)
	if err != nil {
		log.Fatalf("Failed to rename dir: %v", err)
	}
	os.Remove(outputFolder)
	err = os.Rename(tmpPath, outputFolder)
	os.Remove(path)
}
