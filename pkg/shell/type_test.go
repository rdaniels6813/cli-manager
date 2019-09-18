package shell_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/rdaniels6813/cli-manager/pkg/shell"
	"github.com/stretchr/testify/assert"
)

func cleanShellEnv() {
	os.Setenv("ZSH_NAME", "")
	os.Setenv("ZSH", "")
	os.Setenv("BASH", "")
	os.Setenv("PSModulePath", "")
}

func TestGetShellTypeZsh(t *testing.T) {
	shellType := shell.GetShellType(true, false, false, false)
	assert.Equal(t, shell.Zsh, shellType)
}
func TestGetShellTypeBash(t *testing.T) {
	shellType := shell.GetShellType(false, false, true, false)
	assert.Equal(t, shell.Bash, shellType)
}
func TestGetShellTypePowershell(t *testing.T) {
	shellType := shell.GetShellType(false, true, false, false)
	assert.Equal(t, shell.Powershell, shellType)
}
func TestGetShellTypePowershellCore(t *testing.T) {
	shellType := shell.GetShellType(false, false, false, true)
	assert.Equal(t, shell.PowershellCore, shellType)
}

func TestGetShellTypeZshEnv1(t *testing.T) {
	cleanShellEnv()
	os.Setenv("ZSH", "test")
	shellType := shell.GetShellType(false, false, false, false)
	assert.Equal(t, shell.Zsh, shellType)
}
func TestGetShellTypeZshEnv2(t *testing.T) {
	cleanShellEnv()
	os.Setenv("ZSH_NAME", "test")
	shellType := shell.GetShellType(false, false, false, false)
	assert.Equal(t, shell.Zsh, shellType)
}
func TestGetShellTypeBashEnv(t *testing.T) {
	cleanShellEnv()
	os.Setenv("BASH", "test")
	shellType := shell.GetShellType(false, false, false, false)
	assert.Equal(t, shell.Bash, shellType)
}
func TestGetShellTypePowershellEnv(t *testing.T) {
	cleanShellEnv()
	os.Setenv("PSModulePath", "testpath;path2")
	shellType := shell.GetShellType(false, false, false, false)
	assert.Equal(t, shell.Powershell, shellType)
}
func TestGetShellTypePowershellCoreEnv(t *testing.T) {
	cleanShellEnv()
	powershellCoreSubpath := fmt.Sprintf("%spowershell%s", string(os.PathSeparator), string(os.PathSeparator))
	os.Setenv("PSModulePath", powershellCoreSubpath)
	shellType := shell.GetShellType(false, false, false, false)
	assert.Equal(t, shell.PowershellCore, shellType)
}
func TestGetShellTypeUnknown(t *testing.T) {
	cleanShellEnv()
	shellType := shell.GetShellType(false, false, false, false)
	assert.Equal(t, shell.Unknown, shellType)
}
