package fs

import (
	"fmt"
	"log"
	"net/http"
	"scws/common/config"
)

const (
	healthPath = "/_/health"
)

func New(c *config.Config) (*FSStorage, error) {
	fsConfig := config.FsConfig{}
	err := fsConfig.ParseEnv()
	if err != nil {
		log.Println("fs.New", err.Error())
		return nil, err
	}
	s := FSStorage{
		config: &fsConfig,
		index:  c.IndexHtml,
	}
	return &s, nil
}

type FSStorage struct {
	config *config.FsConfig
	index  string
	//scwsConfig *config.Config
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
	if r.URL.Path == healthPath {
		healthHandler(w, r)
		return
	}
	dir := indexDir{
		dir:   s.config.Root,
		index: s.index,
	}
	http.FileServer(dir).ServeHTTP(w, r)
}

func (s *FSStorage) GetName() string {
	return "filesystem"
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
