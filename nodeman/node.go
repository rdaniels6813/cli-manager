package nodeman

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
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
	branchIndex := strings.Index(packageString, "#")
	branch := ""
	if branchIndex == -1 {
		branchIndex = len(packageString)
	} else {
		branch = packageString[branchIndex+1:]
	}
	repoPartialURL := packageString[:branchIndex]
	packageJSONURL := fmt.Sprintf("https://api.github.com/repos/%s/contents/package.json", repoPartialURL)
	if branch != "" {
		packageJSONURL = fmt.Sprintf("%s?ref=%s", packageJSONURL, branch)
	}
	var result NpmViewResponse
	req, err := http.NewRequest("GET", packageJSONURL, nil)
	if err != nil {
		return &result, err
	}
	token := os.Getenv("GH_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("token %s", os.Getenv("GH_TOKEN")))
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &result, err
	}
	if resp.StatusCode != 200 {
		log.Fatalf("Failed to get package.json for project using: %s, please ensure you have set up your github token", packageJSONURL)
	}
	defer resp.Body.Close()
	var content githubContent
	json.NewDecoder(resp.Body).Decode(&content)
	decoded, err := base64.StdEncoding.DecodeString(content.Content)
	if err != nil {
		return &result, nil
	}
	err = json.Unmarshal(decoded, &result)
	if err != nil {
		return &result, nil
	}
	return &result, nil
}

type githubContent struct {
	Content string `json:"content"`
}

// BinPath returns the path to the bin directory for the installed node version
func (n *nodeImpl) BinPath() string {
	return n.getBinPath()
}

func (n *nodeImpl) command(path string, args ...string) error {
	cmd := exec.Command(path)
	var env []string
	for _, val := range os.Environ() {
		name := strings.Split(val, "=")[0]
		if strings.ToLower(name) == "path" {
			val = fmt.Sprintf("%s=%s%s%s", name, filepath.Dir(path), string(os.PathListSeparator), os.Getenv("PATH"))
		}
		env = append(env, val)
	}
	cmd.Env = env
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
