package shell_test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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

func (f *profileFixture) givenFilePathDoesNotExist(fs *mock_afero.MockFs) {
	fs.EXPECT().Stat(f.filePath).Return(nil, os.ErrNotExist)
}
func (f *profileFixture) givenMkdirAllFails(fs *mock_afero.MockFs) {
	fs.EXPECT().MkdirAll(filepath.Dir(f.filePath), os.FileMode(777)).Return(os.ErrPermission)
}
func (f *profileFixture) givenMkdirAllSucceeds(fs *mock_afero.MockFs) {
	fs.EXPECT().MkdirAll(filepath.Dir(f.filePath), os.FileMode(777)).Return(nil)
}
func (f *profileFixture) givenCreateFails(fs *mock_afero.MockFs) {
	fs.EXPECT().Create(f.filePath).Return(nil, os.ErrPermission)
}
func (f *profileFixture) givenCreateSucceeds(fs *mock_afero.MockFs, file afero.File) {
	fs.EXPECT().Create(f.filePath).Return(file, nil)
}
func (f *profileFixture) givenOpenFails(fs *mock_afero.MockFs) {
	fs.EXPECT().Open(f.filePath).Return(nil, os.ErrPermission)
}
func (f *profileFixture) givenOpenFileFails(fs *mock_afero.MockFs) {
	fs.EXPECT().OpenFile(f.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0666)).Return(nil, os.ErrPermission)
}
func (f *profileFixture) givenFilePathExists(fs *mock_afero.MockFs) {
	fs.EXPECT().Stat(f.filePath).Return(nil, nil)
}
func (f *profileFixture) givenOpenSucceeds(fs *mock_afero.MockFs, file afero.File) {
	fs.EXPECT().Open(f.filePath).Return(file, nil)
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

	fixture.givenFilePathDoesNotExist(mockFS)
	fixture.givenMkdirAllFails(mockFS)
	result, err := fixture.helper.WriteProfileSnippet(fixture.snippet, fixture.filePath)

	assert.NotNil(t, err)
	assert.False(t, result)
}

func TestProfileHelperFactory(t *testing.T) {
	helper := shell.NewProfileHelper()

	assert.NotNil(t, helper.FS)
	assert.Equal(t, runtime.GOOS, helper.GOOS)
}

func TestFailToCreateFile(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockFs := mock_afero.NewMockFs(controller)
	fixture := newProfileFixture(mockFs)

	fixture.givenFilePathDoesNotExist(mockFs)
	fixture.givenMkdirAllSucceeds(mockFs)
	fixture.givenCreateFails(mockFs)

	updated, err := fixture.helper.WriteProfileSnippet(fixture.snippet, fixture.filePath)

	assert.False(t, updated)
	assert.Equal(t, os.ErrPermission, err)
}
func TestFailToOpenFileForReading(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockFs := mock_afero.NewMockFs(controller)
	fixture := newProfileFixture(mockFs)

	fixture.givenFilePathExists(mockFs)
	fixture.givenOpenFails(mockFs)

	updated, err := fixture.helper.WriteProfileSnippet(fixture.snippet, fixture.filePath)

	assert.False(t, updated)
	assert.Equal(t, os.ErrPermission, err)
}
func TestFailToOpenFileForWriting(t *testing.T) {
	memfs := afero.NewMemMapFs()
	controller := gomock.NewController(t)
	defer controller.Finish()
	mockFs := mock_afero.NewMockFs(controller)
	fixture := newProfileFixture(mockFs)

	fixture.givenFilePathExists(mockFs)
	testFile, err := memfs.Create("testfile")
	assert.Nil(t, err)
	fixture.givenOpenSucceeds(mockFs, testFile)
	fixture.givenOpenFileFails(mockFs)

	updated, err := fixture.helper.WriteProfileSnippet(fixture.snippet, fixture.filePath)

	assert.False(t, updated)
	assert.Equal(t, os.ErrPermission, err)
}
