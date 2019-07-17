package nodeman

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

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

// CLIApp config for an installed CLI app
type CLIApp struct {
	Bin  string `json:"bin"`
	Path string `json:"path"`
}

// GetInstalledApps list all of the installed apps
func (m *Manager) GetInstalledApps() []string {
	configPath := m.getConfigPath()
	apps := loadConfig(configPath)
	var result []string
	for app := range apps {
		result = append(result, app)
	}
	return result
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

func (m *Manager) getCliManagerFolder() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("There was an error determining the user home directory")
	}
	cliManagerDir := filepath.Join(homeDir, ".cli-manager")
	if _, err := m.os.Stat(cliManagerDir); os.IsNotExist(err) {
		m.os.MkdirAll(cliManagerDir, 0700)
	}
	return cliManagerDir
}

// GetCommandPath gets the path to the installed command
func (m *Manager) GetCommandPath(bin string) (string, error) {
	cliManagerDir := m.getCliManagerFolder()
	installedAppsJSON := filepath.Join(cliManagerDir, "installed.json")
	config := loadConfig(installedAppsJSON)

	for k, v := range config {
		if k == bin {
			if runtime.GOOS == "windows" {
				return filepath.Join(v.Path, fmt.Sprintf("%s.cmd", k)), nil
			}
			return filepath.Join(v.Path, fmt.Sprintf("%s", k)), nil
		}
	}
	return "", fmt.Errorf("%s is not installed", bin)
}

// ConfigureNodeOnCommand configures a command to use the appropriate PATH
// var for the node version required by a command
func (m *Manager) ConfigureNodeOnCommand(command string, cmd *exec.Cmd) {
	nodeBinPath, err := m.GetCommandNodeBinPath(command)
	if err != nil {
		log.Fatal(err)
	}
	var env []string
	for _, val := range os.Environ() {
		name := strings.Split(val, "=")[0]
		if strings.ToLower(name) == "path" {
			val = fmt.Sprintf("%s=%s%s%s", name, nodeBinPath, string(os.PathListSeparator), os.Getenv("PATH"))
		}
		env = append(env, val)
	}
	cmd.Env = env
}

// GetCommandNodeBinPath returns the bin folder path to the node installation being used for the command
func (m *Manager) GetCommandNodeBinPath(bin string) (string, error) {
	cliManagerDir := m.getCliManagerFolder()
	installedAppsJSON := filepath.Join(cliManagerDir, "installed.json")
	config := loadConfig(installedAppsJSON)

	for k, v := range config {
		if k == bin {
			return v.Path, nil
		}
	}
	return "", fmt.Errorf("%s is not installed", bin)
}

func (m *Manager) getConfigPath() string {
	cliManagerDir := m.getCliManagerFolder()
	return filepath.Join(cliManagerDir, "installed.json")
}

// MarkInstalled marks the current binaries installed with their paths
func (m *Manager) MarkInstalled(bins map[string]string, binpath string) error {
	installedAppsJSON := m.getConfigPath()
	config := loadConfig(installedAppsJSON)
	for k := range bins {
		config[k] = &CLIApp{Path: binpath, Bin: k}
	}
	return saveConfig(installedAppsJSON, config)
}

func saveConfig(path string, config map[string]*CLIApp) error {
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return err
	}
	err = json.NewEncoder(f).Encode(config)
	if err != nil {
		return err
	}
	return nil
}

func loadConfig(path string) map[string]*CLIApp {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return make(map[string]*CLIApp)
	}
	var result map[string]*CLIApp
	err = json.NewDecoder(f).Decode(&result)
	if err != nil {
		return make(map[string]*CLIApp)
	}
	return result
}

func (m *Manager) getNodeBaseFolder() string {
	cliManagerDir := m.getCliManagerFolder()
	nodeFolder := filepath.Join(cliManagerDir, "node")
	if _, err := m.os.Stat(nodeFolder); os.IsNotExist(err) {
		m.os.MkdirAll(nodeFolder, 0700)
	}
	return nodeFolder
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
