package nodeman

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rdaniels6813/cli-manager/internal/util"
	"github.com/spf13/afero"

	"github.com/mholt/archiver"
)

const windows = "windows"

// Manager struct for managing node binaries leveraged by the cli manager
// Use NewManager to create an instance of this object
type Manager struct {
	os afero.Fs
}

// NewManager constructor for default manager with the specified node version
func NewManager(os afero.Fs) *Manager {
	return &Manager{os: os}
}

// CLIApp config for an installed CLI app
type CLIApp struct {
	App         string `json:"app"`
	Bin         string `json:"bin"`
	Path        string `json:"path"`
	InstallName string `json:"install_name"`
}

// GetInstalledExecutables list all of the installed app executables
func (m *Manager) GetInstalledExecutables() []string {
	configPath := m.getConfigPath()
	apps := loadConfig(configPath)
	result := make([]string, 0, len(apps))
	for app := range apps {
		result = append(result, app)
	}
	return result
}

// GetNode Ensures the node binary is downloaded & extracted and returns the path to the bin folder
func (m *Manager) GetNode(version string) (Node, error) {
	destinationPath := m.getNodeOutputFolder(version)
	if _, err := m.os.Stat(destinationPath); os.IsNotExist(err) {
		archivePath, err := m.downloadNodeArchive(version)
		if err != nil {
			return nil, err
		}
		m.unpackNodeArchive(archivePath, version)
	}
	return newNode(destinationPath), nil
}

// GetNodeByPath uses the binpath to set up a node helper
func (m *Manager) GetNodeByPath(path string) Node {
	return newNode(path)
}

// GetCLIApp gets the configuration for the cli app by name
func (m *Manager) GetCLIApp(appName string) (*CLIApp, error) {
	configPath := m.getConfigPath()
	apps := loadConfig(configPath)
	for _, app := range apps {
		if app.App == appName {
			return app, nil
		}
	}
	for _, app := range apps {
		if app.InstallName == appName {
			return app, nil
		}
	}
	return nil, fmt.Errorf("App is not installed: %s", appName)
}

func (m *Manager) getNodeURL(version string) string {
	extension := ".tar.xz"
	os := runtime.GOOS
	arch := runtime.GOARCH
	if os == windows {
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

// GetCommandPath gets the path to the installed command
func (m *Manager) GetCommandPath(bin string) (string, error) {
	cliManagerDir := util.GetCliManagerFolder(m.os)
	installedAppsJSON := filepath.Join(cliManagerDir, "installed.json")
	config := loadConfig(installedAppsJSON)

	for k, v := range config {
		if k == bin {
			if runtime.GOOS == windows {
				return filepath.Join(v.Path, fmt.Sprintf("%s.cmd", k)), nil
			}
			return filepath.Join(v.Path, k), nil
		}
	}
	return "", fmt.Errorf("%s is not installed", bin)
}

// ConfigureNodeOnCommand configures a command to use the appropriate PATH
// var for the node version required by a command
func (m *Manager) ConfigureNodeOnCommand(command string, cmd *exec.Cmd) {
	nodeBinPath, err := m.GetCommandNodeBinPath(command)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	env := make([]string, 0, len(os.Environ()))
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
	cliManagerDir := util.GetCliManagerFolder(m.os)
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
	cliManagerDir := util.GetCliManagerFolder(m.os)
	return filepath.Join(cliManagerDir, "installed.json")
}

// MarkUninstalled marks the current binaries installed with their paths
func (m *Manager) MarkUninstalled(appName string) error {
	installedAppsJSON := m.getConfigPath()
	config := loadConfig(installedAppsJSON)
	for bin, app := range config {
		if app.App == appName {
			delete(config, bin)
		}
	}
	return saveConfig(installedAppsJSON, config)
}

// MarkInstalled marks the current binaries installed with their paths
func (m *Manager) MarkInstalled(app string, bins map[string]string, binpath string, installName string) error {
	installedAppsJSON := m.getConfigPath()
	config := loadConfig(installedAppsJSON)
	for k := range bins {
		config[k] = &CLIApp{Path: binpath, Bin: k, App: app, InstallName: installName}
	}
	return saveConfig(installedAppsJSON, config)
}

func saveConfig(path string, config map[string]*CLIApp) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(config)
	if err != nil {
		return err
	}
	return nil
}

func loadConfig(path string) map[string]*CLIApp {
	f, err := os.Open(path)
	if err != nil {
		return make(map[string]*CLIApp)
	}
	defer f.Close()
	var result map[string]*CLIApp
	err = json.NewDecoder(f).Decode(&result)
	if err != nil {
		return make(map[string]*CLIApp)
	}
	return result
}

func (m *Manager) getNodeBaseFolder() string {
	cliManagerDir := util.GetCliManagerFolder(m.os)
	nodeFolder := filepath.Join(cliManagerDir, "node")
	if _, err := m.os.Stat(nodeFolder); os.IsNotExist(err) {
		err = m.os.MkdirAll(nodeFolder, 0700)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nodeFolder
}

func (m *Manager) downloadNodeArchive(version string) (string, error) {
	nodeBaseFolder := m.getNodeBaseFolder()
	// Create the file
	nodeBinaryPath := filepath.Join(nodeBaseFolder, filepath.Base(m.getNodeURL(version)))
	if _, err := m.os.Stat(nodeBinaryPath); os.IsNotExist(err) {
		out, err := os.Create(nodeBinaryPath)
		if err != nil {
			return "", fmt.Errorf("Failed to create destination file for node binary: %w", err)
		}
		defer out.Close()

		// Get the data
		resp, err := http.Get(m.getNodeURL(version))
		if err != nil {
			return "", fmt.Errorf("Failed download node binary: %w", err)
		}
		defer resp.Body.Close()

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return "", fmt.Errorf("Failed write node binary: %w", err)
		}
	}
	return nodeBinaryPath, nil
}
func (m *Manager) getNodeOutputFolder(version string) string {
	nodeFolder := m.getNodeBaseFolder()
	return filepath.Join(nodeFolder, version)
}

func (m *Manager) unpackNodeArchive(path string, version string) {
	outputFolder := m.getNodeOutputFolder(version)
	err := archiver.Unarchive(path, outputFolder)
	if err != nil {
		fmt.Printf("Failed to unarchive: %s", err)
		os.Exit(1)
	}
	dirs, err := os.ReadDir(outputFolder)
	if err != nil {
		fmt.Printf("Failed to read dir: %v", err)
		os.Exit(1)
	}
	dirName := filepath.Join(outputFolder, dirs[0].Name())
	tmpPath, err := os.MkdirTemp(m.getNodeBaseFolder(), version)
	os.Remove(tmpPath)
	if err != nil {
		fmt.Printf("Failed to get temp dir: %v", err)
		os.Exit(1)
	}
	err = os.Rename(dirName, tmpPath)
	if err != nil {
		fmt.Printf("Failed to rename dir: %v", err)
		os.Exit(1)
	}
	os.Remove(outputFolder)
	err = os.Rename(tmpPath, outputFolder)
	if err != nil {
		fmt.Printf("Failed to rename archive paths: %s", err)
		os.Exit(1)
	}
	os.Remove(path)
}
