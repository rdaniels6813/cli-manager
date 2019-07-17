package nodeman

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Node helper for wrapping node commands for a node installation
type Node interface {
	Node(args ...string) error
	Npm(args ...string) error
	NpmView(packageString string) (*NpmViewResponse, error)
	BinPath() string
}

type nodeImpl struct {
	nodePath string
}

// Node execute a command using the node binary with the following arguments `node args[0] args[1] ...`
func (n *nodeImpl) Node(args ...string) error {
	return n.command(n.getNodePath(), args...)
}

// Npm execute a command using npm with the following arguments `npm args[0] args[1] ...`
func (n *nodeImpl) Npm(args ...string) error {
	return n.command(n.getNpmPath(), args...)
}

// NpmViewResponse response from npm view command
type NpmViewResponse struct {
	Engines map[string]string `json:"engines"`
	Bin     map[string]string `json:"bin"`
}

// Npm execute a command using npm with the following arguments `npm args[0] args[1] ...` returning results as JSON
func (n *nodeImpl) NpmView(packageString string) (*NpmViewResponse, error) {
	cmd := exec.Command(n.getNpmPath())
	cmd.Args = append(cmd.Args, "view", packageString)
	cmd.Args = append(cmd.Args, "--json")
	output, err := cmd.Output()
	if err != nil {
		return n.githubPackageJSON(packageString)
	}
	var response NpmViewResponse
	err = json.Unmarshal(output, &response)
	return &response, err
}

func (n *nodeImpl) githubPackageJSON(packageString string) (*NpmViewResponse, error) {
	repoPartialURL := strings.ReplaceAll(packageString, "#", "/")
	packageJSONURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/package.json", repoPartialURL)
	var result NpmViewResponse
	resp, err := http.Get(packageJSONURL)
	if err != nil {
		return &result, err
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&result)
	return &result, nil
}

// BinPath returns the path to the bin directory for the installed node version
func (n *nodeImpl) BinPath() string {
	return n.getBinPath()
}

func (n *nodeImpl) command(path string, args ...string) error {
	cmd := exec.Command(path)
	cmd.Env = append(os.Environ(), fmt.Sprintf("npm_config_prefix=\"%s\"", strings.ReplaceAll(n.nodePath, "\\", "/")))
	cmd.Args = append(cmd.Args, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (n *nodeImpl) getBinPath() string {
	if runtime.GOOS == "windows" {
		return n.nodePath
	}
	return filepath.Join(n.nodePath, "bin/")
}

func (n *nodeImpl) getNodePath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(n.nodePath, "node.exe")
	}
	return filepath.Join(n.nodePath, "bin/node")
}

func (n *nodeImpl) getNpmPath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(n.nodePath, "npm.cmd")
	}
	return filepath.Join(n.nodePath, "bin/npm")
}

// NewNode create a new NodeHelper based on the path to the bin folder
func newNode(nodePath string) Node {
	return &nodeImpl{nodePath: nodePath}
}
