package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spf13/afero"
)

func TestGetNodeBaseFolder(t *testing.T) {
	fs := afero.NewMemMapFs()
	app := &App{os: fs}
	app.getNodeBaseFolder()
	homedir, _ := os.UserHomeDir()
	stat, err := fs.Stat(filepath.Join(homedir, ".cli-manager", "node"))
	assert.Nil(t, err)
	assert.True(t, stat.IsDir())
}
