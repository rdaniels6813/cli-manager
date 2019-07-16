package nodeman

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Node helper for wrapping node commands for a node installation
type Node interface {
	Node(args ...string) error
	Npm(args ...string) error
	NpmView(packageString string) (*NpmViewResponse, error)
}

type nodeImpl struct {
	nodePath string
}

// Node execute a command using the node binary with the following arguments `node args[0] args[1] ...`
func (n *nodeImpl) Node(args ...string) error {
	return command(n.getNodePath(), args...)
}

// Npm execute a command using npm with the following arguments `npm args[0] args[1] ...`
func (n *nodeImpl) Npm(args ...string) error {
	command(n.getNpmPath(), "config", "set", "prefix", n.nodePath)
	return command(n.getNpmPath(), args...)
}

// NpmViewResponse response from npm view command
type NpmViewResponse struct {
	Engines map[string]string `json:"engines"`
}

// Npm execute a command using npm with the following arguments `npm args[0] args[1] ...` returning results as JSON
func (n *nodeImpl) NpmView(packageString string) (*NpmViewResponse, error) {
	cmd := exec.Command(n.getNpmPath())
	cmd.Args = append(cmd.Args, "view", packageString)
	cmd.Args = append(cmd.Args, "--json")
	output, err := cmd.Output()
	var response NpmViewResponse
	err = json.Unmarshal(output, &response)
	return &response, err
}

func command(path string, args ...string) error {
	cmd := exec.Command(path)
	cmd.Args = append(cmd.Args, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
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
