package fs

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path"
	"scws/config"
)

func New(index string) (*FSStorage, error) {
	fsConfig := config.FsConfig{}
	err := fsConfig.ParseEnv()
	if err != nil {
		log.Println("fs.New", err.Error())
		return nil, err
	}
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
	v := http.Dir(d.dir)
	file, err := v.Open(name)
	if err != nil {
		file, err = v.Open(d.index)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
	return file, nil
}

func (s *FSStorage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dir := indexDir{
		dir:   s.config.Root,
		index: s.index,
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
