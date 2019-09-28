package aliases

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/rdaniels6813/cli-manager/pkg/nodeman"
	"github.com/rdaniels6813/cli-manager/pkg/shell"
)

func NewZshGenerator() Generator {
	return &ZshGenerator{}
}

type ZshGenerator struct {
	NodeManager nodeman.Manager
}

func (g *ZshGenerator) Generate() string {
	apps := g.NodeManager.GetInstalledExecutables()
	result := ""
	for _, app := range apps {
		result += fmt.Sprintf("alias %s='cli-manager run %s'\n", app, app)
	}
	return result
}

const ZshAliasesSnippet = "source <(cli-manager aliases -g -z)\n"

func (g *ZshGenerator) Install() error {
	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	scriptPath := filepath.Join(dir, ".zshrc")
	profileHelper := shell.ProfileHelper{}
	wrote, err := profileHelper.WriteProfileSnippet(ZshAliasesSnippet, scriptPath)
	if err != nil {
		log.Fatal(err)
	}
	if wrote {
		fmt.Printf("Wrote aliases script to: %s\n", scriptPath)
	} else {
		fmt.Printf("Aliases already installed in: %s\n", scriptPath)
	}
	return nil
}
