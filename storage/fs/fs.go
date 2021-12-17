package fs

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"scws/config"
	"strings"
)

func New(index string) (*FSStorage, error) {
	fsConfig := config.FsConfig{}
	fsConfig.ParseEnv()
	s := FSStorage{
		config: &fsConfig,
		index:  index,
	}
	return &s, nil
}

type FSStorage struct {
	config *config.FsConfig
	index  string
}

type indexDir struct {
	dir   string
	index string
}

func (d indexDir) Open(name string) (http.File, error) {
	var (
		file http.File
		err  error
	)
	v := http.Dir(d.dir)
	file, err = v.Open(name)

	if errors.Is(err, fs.ErrNotExist) {
		file, err = v.Open(d.index)
		if err != nil {
			return nil, fmt.Errorf("file can't be opened %v", err)
		}
		return file, nil
	}
	if err != nil {
		return nil, fmt.Errorf("file can't be opened %v", err)
	}
	return file, nil
}

func (s *FSStorage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dir := indexDir{
		dir:   s.config.Root,
		index: s.index,
	}
	if strings.HasSuffix(r.URL.Path, "/") {
		r.URL.Path = "/"
	}
	http.FileServer(dir).ServeHTTP(w, r)
}

func (s *FSStorage) GetName() string {
	return "filesystem"
}

func (s *FSStorage) indexPath() string {
	return path.Join(s.config.Root, s.index)
}

func checkFile(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return err
	} else {
		return nil
	}
}

func (s *FSStorage) HealthProbe() error {
	return checkFile(s.indexPath())
}
