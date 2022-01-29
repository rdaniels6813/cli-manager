package store

import (
	"encoding/json"
	"os"
	"path"

	"github.com/rdaniels6813/cli-manager/internal/util"
	"github.com/spf13/afero"
)

type Store interface {
	Get(hostname string, key string) (string, error)
	Set(hostname string, key string, value string) error
}

type FSStore struct {
	path string
}

func GetDefaultStore() Store {
	return &FSStore{
		path: path.Join(util.GetCliManagerFolder(afero.OsFs{}), "auth.json"),
	}
}

func (f *FSStore) Get(hostname string, key string) (string, error) {
	cfg, err := os.Open(f.path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	defer cfg.Close()
	d := json.NewDecoder(cfg)
	values := map[string]map[string]string{}
	err = d.Decode(&values)
	if err != nil {
		return "", err
	}
	return values[hostname][key], nil
}

func (f *FSStore) Set(hostname string, key string, value string) error {
	cfg, err := os.OpenFile(f.path, os.O_RDWR, 0600)
	values := map[string]map[string]string{}
	if os.IsNotExist(err) {
		cfg, err = os.OpenFile(f.path, os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			return err
		}
	} else {
		if err != nil {
			return err
		}
		d := json.NewDecoder(cfg)
		err = d.Decode(&values)
		if err != nil && err.Error() != "EOF" {
			return err
		}
	}
	defer cfg.Close()
	if values[hostname] == nil {
		values[hostname] = map[string]string{}
	}
	values[hostname][key] = value
	_, err = cfg.Seek(0, 0)
	if err != nil {
		return err
	}
	e := json.NewEncoder(cfg)
	return e.Encode(&values)
}
