package shell_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rdaniels6813/cli-manager/pkg/mock_afero"
	"github.com/rdaniels6813/cli-manager/pkg/shell"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

type profileFixture struct {
	fs       afero.Fs
	helper   *shell.ProfileHelper
	snippet  string
	filePath string
}

func newProfileFixture(fs afero.Fs) *profileFixture {
	helper := &shell.ProfileHelper{
		FS: fs,
	}
	return &profileFixture{
		fs:       fs,
		helper:   helper,
		filePath: "/fake/path/.bashrc",
		snippet:  "test snippet",
	}
}

func TestWriteProfileSnippetNewFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	fixture := newProfileFixture(fs)
	result, err := fixture.helper.WriteProfileSnippet(fixture.snippet, fixture.filePath)

	assert.Nil(t, err)
	assert.True(t, result)

	profileBytes, err := afero.ReadFile(fixture.fs, fixture.filePath)
	assert.Nil(t, err)
	profileString := string(profileBytes)
	assert.Contains(t, profileString, fixture.snippet)
}

func TestWriteProfileSnippetExistingFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	fixture := newProfileFixture(fs)
	existingProfile := "# this is my profile\necho Welcome!\n"

	err := afero.WriteFile(fixture.fs, fixture.filePath, []byte(existingProfile), 0666)
	assert.Nil(t, err)

	result, err := fixture.helper.WriteProfileSnippet(fixture.snippet, fixture.filePath)

	assert.Nil(t, err)
	assert.True(t, result)

	profileBytes, err := afero.ReadFile(fixture.fs, fixture.filePath)
	assert.Nil(t, err)
	profileString := string(profileBytes)
	assert.Contains(t, profileString, fixture.snippet)
	assert.Contains(t, profileString, existingProfile)
}

func TestWriteProfileSnippetExistingFileAlreadyHasSnippet(t *testing.T) {
	fs := afero.NewMemMapFs()
	fixture := newProfileFixture(fs)
	existingProfile := fmt.Sprintf("# this is my profile\necho Welcome!\n%s\n", fixture.snippet)

	err := afero.WriteFile(fixture.fs, fixture.filePath, []byte(existingProfile), 0666)
	assert.Nil(t, err)

	result, err := fixture.helper.WriteProfileSnippet(fixture.snippet, fixture.filePath)

	assert.Nil(t, err)
	assert.False(t, result)
}

func TestWriteProfileSnippetFailsToMakeDirectory(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockFS := mock_afero.NewMockFs(controller)
	fixture := newProfileFixture(mockFS)

	mockFS.EXPECT().Stat(fixture.filePath).Times(1).Return(nil, os.ErrNotExist)
	mockFS.EXPECT().MkdirAll(filepath.Dir(fixture.filePath), os.FileMode(777)).Times(1).Return(os.ErrPermission)
	result, err := fixture.helper.WriteProfileSnippet(fixture.snippet, fixture.filePath)

	assert.NotNil(t, err)
	assert.False(t, result)
}
