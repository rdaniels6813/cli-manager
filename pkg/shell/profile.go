package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/afero"
)

type ProfileHelper struct {
	FS   afero.Fs
	GOOS string
}

func NewProfileHelper() *ProfileHelper {
	return &ProfileHelper{
		FS:   afero.OsFs{},
		GOOS: runtime.GOOS,
	}
}

func (p *ProfileHelper) WriteProfileSnippet(snippet, path string) (bool, error) {
	if _, err := p.FS.Stat(path); os.IsNotExist(err) {
		err := p.FS.MkdirAll(filepath.Dir(path), 777)
		if err != nil {
			return false, err
		}
		f, err := p.FS.Create(path)
		if err != nil {
			return false, err
		}
		defer f.Close()
		_, err = f.WriteString(snippet)
		return true, err
	}
	profileBytes, err := afero.ReadFile(p.FS, path)
	if err != nil {
		return false, err
	}
	text := string(profileBytes)
	if strings.Contains(text, snippet) {
		return false, nil
	}
	f, err := p.FS.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return false, err
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("%s\n", snippet))
	return true, err
}
