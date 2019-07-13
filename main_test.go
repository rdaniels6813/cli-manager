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

func TestGetNodeURL(t *testing.T) {
	version := "10.16.0"
	app := NewApp()
	t.Run("darwin-amd64", func(t *testing.T) {
		url := app.getNodeURL(version, "darwin", "amd64")
		assert.Equal(t, "https://nodejs.org/dist/v10.16.0/node-v10.16.0-darwin-x64.tar.gz", url)
	})
	t.Run("windows-amd64", func(t *testing.T) {
		url := app.getNodeURL(version, "windows", "amd64")
		assert.Equal(t, "https://nodejs.org/dist/v10.16.0/node-v10.16.0-win-x64.zip", url)
	})
	t.Run("linux-amd64", func(t *testing.T) {
		url := app.getNodeURL(version, "linux", "amd64")
		assert.Equal(t, "https://nodejs.org/dist/v10.16.0/node-v10.16.0-linux-x64.tar.xz", url)
	})
}

func TestUnpackNodeBinary(t *testing.T) {
	app := NewApp()
	t.Run("darwin-amd64", func(t *testing.T) {
		path := app.unpackNodeBinary("/tmp/node-v10.16.0-darwin-x64.tar.gz")
	})
	t.Run("windows-amd64", func(t *testing.T) {
		path := app.unpackNodeBinary("/tmp/node-v10.16.0-win-x64.zip")
	})
	t.Run("linux-amd64", func(t *testing.T) {
		path := app.unpackNodeBinary("/tmp/node-v10.16.0-linux-x64.tar.xz")
	})
}
