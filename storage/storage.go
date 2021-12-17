package storage

import (
	"fmt"
	"net/http"
	"scws/config"
	"scws/storage/fs"
	"scws/storage/s3"

	"github.com/pkg/errors"
)

const (
	FSStorage = "filesystem"
	S3        = "s3"
)

type StorageHandler interface {
	Handler() http.Handler
	HealthProbe() http.Handler
}

type IStorage interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	GetName() string
	HealthProbe() error
}

type Storage struct {
	storage IStorage
}

func New(c *config.Config) (*Storage, error) {
	var err error
	s := Storage{}
	switch c.Storage {
	case FSStorage:
		s.storage, err = fs.New(c.IndexHtml)
	case S3:
		s.storage, err = s3.New(c.IsVaultEnabled(), c.VaultPaths)
	}
	if err != nil {
		return nil, fmt.Errorf("couldn't initialize storage: %v", err)
	}
	if s.storage == nil {
		return nil, errors.New("storage.New has unexpected error")
	}
	return &s, nil
}

func (s *Storage) Handler() http.Handler {
	return s
}

func (s *Storage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.storage.ServeHTTP(w, r)
}

func (s *Storage) HealthProbe() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := s.HealthProbe(); err != nil {
			fmt.Fprint(w, "healthy")
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	})
}
