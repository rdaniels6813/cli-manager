package aliases

import (
	"fmt"
	"log"

	"github.com/rdaniels6813/cli-manager/pkg/nodeman"
	"github.com/rdaniels6813/cli-manager/pkg/shell"
)

func NewPowershellGenerator(powershellCore bool) Generator {
	return &PowershellGenerator{PowershellCore: powershellCore}
}

type PowershellGenerator struct {
	NodeManager    nodeman.Manager
	PowershellCore bool
}

func (g *PowershellGenerator) Generate() string {
	apps := g.NodeManager.GetInstalledExecutables()
	result := ""
	for _, app := range apps {
		result += fmt.Sprintf("function %s { cli-manager.exe run %s @args }", app, app)
	}
	return result
}

func (g *PowershellGenerator) Install() error {
	scriptPath := shell.GetPowershellProfilePath(g.PowershellCore)
	wrote, err := shell.WriteProfileSnippet(PowershellAliasesSnippet, scriptPath)
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

const PowershellAliasesSnippet = "Invoke-Expression $($(cli-manager.exe aliases -g -p) -join \"`n\")\n"
