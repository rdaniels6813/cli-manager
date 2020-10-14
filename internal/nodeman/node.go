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

const WINDOWS = "windows"

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
	Name    string            `json:"name"`
	Engines map[string]string `json:"engines"`
	Bin     interface{}       `json:"bin"`
}

func (n *NpmViewResponse) GetBins() map[string]string {
	if binString, ok := n.Bin.(string); ok {
		result := make(map[string]string, 1)
		result[filepath.Base(binString)] = binString
		return result
	}
	if binMap, ok := n.Bin.(map[string]string); ok {
		return binMap
	}
	if interfaceMap, ok := n.Bin.(map[string]interface{}); ok {
		result := make(map[string]string, len(interfaceMap))
		for k, v := range interfaceMap {
			result[k] = v.(string)
		}
		return result
	}
	return nil
}

// Npm execute a command using npm with the following arguments `npm args[0] args[1] ...` returning results as JSON
func (n *nodeImpl) NpmView(packageString string) (*NpmViewResponse, error) {
	// This command runs a local node version based on a calculated path.
	cmd := exec.Command(n.getNpmPath()) //nolint:gosec
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
		log.Fatalf("Failed to get package.json for project using: %s, please ensure you have set up your github token",
			packageJSONURL)
	}
	defer resp.Body.Close()
	var content githubContent
	err = json.NewDecoder(resp.Body).Decode(&content)
	if err != nil {
		return nil, err
	}
	decoded, err := base64.StdEncoding.DecodeString(content.Content)
	if err != nil {
		return nil, err
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
	env := make([]string, 0, len(os.Environ()))
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
	if runtime.GOOS == WINDOWS {
		return n.nodePath
	}
	return filepath.Join(n.nodePath, "bin/")
}

func (n *nodeImpl) getNodePath() string {
	if runtime.GOOS == WINDOWS {
		return filepath.Join(n.nodePath, "node.exe")
	}
	if strings.HasSuffix(n.nodePath, "/bin") {
		return filepath.Join(n.nodePath, "node")
	}
	return filepath.Join(n.nodePath, "bin/node")
}

func (n *nodeImpl) getNpmPath() string {
	if runtime.GOOS == WINDOWS {
		return filepath.Join(n.nodePath, "npm.cmd")
	}
	if strings.HasSuffix(n.nodePath, "/bin") {
		return filepath.Join(n.nodePath, "npm")
	}
	return filepath.Join(n.nodePath, "bin/npm")
}

// NewNode create a new NodeHelper based on the path to the bin folder
func newNode(nodePath string) Node {
	return &nodeImpl{nodePath: nodePath}
}
