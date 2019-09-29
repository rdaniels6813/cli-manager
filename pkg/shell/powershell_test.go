package shell_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rdaniels6813/cli-manager/pkg/shell"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type powershellFixture struct {
	fs            afero.Fs
	profileHelper *shell.ProfileHelper
}

func newPowershellFixture(fs afero.Fs) *powershellFixture {
	return &powershellFixture{
		fs: fs,
		profileHelper: &shell.ProfileHelper{
			FS: fs,
		},
	}
}

func TestGetPowershellCoreProfilePathDarwin(t *testing.T) {
	fs := afero.NewMemMapFs()
	fixture := newPowershellFixture(fs)
	fixture.profileHelper.GOOS = "darwin"

	homeDir, _ := os.UserHomeDir()
	expectedPath := filepath.Join(homeDir, ".config", "powershell", "Microsoft.PowerShell_profile.ps1")
	err := fs.MkdirAll(filepath.Dir(expectedPath), 777)
	assert.Nil(t, err)

	pathResult, err := fixture.profileHelper.GetPowershellProfilePath(true)
	assert.Nil(t, err)
	assert.Equal(t, expectedPath, pathResult)
}

func TestGetPowershellCoreProfilePathWindows(t *testing.T) {
	fs := afero.NewMemMapFs()
	fixture := newPowershellFixture(fs)
	fixture.profileHelper.GOOS = "windows"

	homeDir, _ := os.UserHomeDir()
	expectedPath := filepath.Join(homeDir, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1")
	err := fs.MkdirAll(filepath.Dir(expectedPath), 777)
	assert.Nil(t, err)

	pathResult, err := fixture.profileHelper.GetPowershellProfilePath(true)

	assert.Nil(t, err)
	assert.Equal(t, expectedPath, pathResult)
}

func TestGetPowershellCoreProfilePathLinux(t *testing.T) {
	fs := afero.NewMemMapFs()
	fixture := newPowershellFixture(fs)
	fixture.profileHelper.GOOS = "linux"

	homeDir, _ := os.UserHomeDir()
	expectedPath := filepath.Join(homeDir, ".config", "powershell", "Microsoft.PowerShell_profile.ps1")
	err := fs.MkdirAll(filepath.Dir(expectedPath), 777)
	assert.Nil(t, err)

	pathResult, err := fixture.profileHelper.GetPowershellProfilePath(true)

	assert.Nil(t, err)
	assert.Equal(t, expectedPath, pathResult)
}
