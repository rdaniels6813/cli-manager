package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/go-github/v33/github"
	"github.com/rdaniels6813/cli-manager/internal/version"
	"github.com/spf13/cobra"
)

const (
	NAME_TEMPLATE = `cli-manager-%s-%s%s`
	OWNER         = "rdaniels6813"
	REPO          = "cli-manager"
)

// upgradeCmd represents the command to update cli-manager
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade cli-manager to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		latest, err := getLatestReleaseBinary()
		if err != nil {
			return err
		}
		err = latest.Chmod(0755)
		if err != nil {
			latest.Close()
			return fmt.Errorf("Failed to update binary permissions; %w", err)
		}
		latest.Close()
		newExecName := latest.Name()
		err = testNewExecutable(newExecName)
		if err != nil {
			return err
		}
		filePath := newExecName
		originPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("Failed get path for original executable; %w", err)
		}
		resolvedPath, err := filepath.EvalSymlinks(originPath)
		if err != nil {
			return fmt.Errorf("Failed to resolve symlinks for original executable; %w", err)
		}
		err = os.Rename(resolvedPath, resolvedPath+".bak")
		if err != nil {
			return fmt.Errorf("Failed to move original executable; %w", err)
		}
		err = os.Rename(filePath, resolvedPath)
		if err != nil {
			return fmt.Errorf("Failed to move new executable; %w", err)
		}
		err = os.Remove(resolvedPath + ".bak")
		if err != nil {
			return fmt.Errorf("Failed to remove old executable; %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

func testNewExecutable(path string) error {
	output, err := exec.Command(path, "--version").CombinedOutput()
	if err != nil {
		return err
	}
	if !strings.Contains(string(output), "cli-manager version ") {
		return fmt.Errorf("New cli version test failed, try again later or download an update manually")
	}
	return nil
}

func getLatestReleaseBinary() (*os.File, error) {
	client := github.NewClient(nil)
	latestRelease, _, err := client.Repositories.GetLatestRelease(context.Background(), OWNER, REPO)
	if err != nil {
		return nil, err
	}
	if *latestRelease.TagName == version.GetModuleVersion() {
		return nil, fmt.Errorf("Already on the latest version")
	}
	extension := ""
	if runtime.GOOS == "windows" {
		extension = ".exe"
	}
	desiredBinary := fmt.Sprintf(NAME_TEMPLATE, runtime.GOOS, runtime.GOARCH, extension)
	for _, asset := range latestRelease.Assets {
		if asset != nil && asset.Name != nil && *asset.Name == desiredBinary {
			rc, _, err := client.Repositories.DownloadReleaseAsset(context.Background(), OWNER, REPO, *asset.ID, http.DefaultClient)
			if err != nil {
				return nil, fmt.Errorf("Failed getting asset to download; %w", err)
			}
			defer rc.Close()
			p, err := os.Executable()
			if err != nil {
				return nil, fmt.Errorf("Failed getting executable path; %w", err)
			}
			f, err := os.CreateTemp(path.Dir(p), "cli-manager")
			if err != nil {
				return nil, fmt.Errorf("Failed creating temp file; %w", err)
			}
			_, err = io.Copy(f, rc)
			if err != nil {
				f.Close()
				return nil, fmt.Errorf("Failed copying dowload to temp file; %w", err)
			}
			return f, nil
		}
	}
	return nil, fmt.Errorf("Couldn't find newer version to download")
}
