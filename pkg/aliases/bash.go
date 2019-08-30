package aliases

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rdaniels6813/cli-manager/pkg/nodeman"
	"github.com/rdaniels6813/cli-manager/pkg/shell"
)

func NewBashGenerator() Generator {
	return &BashGenerator{}
}

type BashGenerator struct {
	NodeManager nodeman.Manager
}

func (g *BashGenerator) Generate() string {
	apps := g.NodeManager.GetInstalledExecutables()
	result := ""
	for _, app := range apps {
		result += fmt.Sprintf("alias %s='cli-manager run %s'\n", app, app)
	}
	return result
}

func (g *BashGenerator) Install() error {

	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	scriptPath := filepath.Join(dir, ".bashrc")
	wrote, err := shell.WriteProfileSnippet(BashAliasesSnippet, scriptPath)
	if err != nil {
		return err
	}
	if wrote {
		fmt.Printf("Wrote aliases script to: %s\n", scriptPath)
	} else {
		fmt.Printf("Aliases already installed in: %s\n", scriptPath)
	}
	return nil
}

const BashAliasesSnippet = "source <(cli-manager aliases -g -b)\n"
