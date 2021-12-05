package storage

import (
	"log"
	"net/http"
	"scws/config"
	"scws/storage/fs"
	"scws/storage/s3"
)

const (
	FSStorage = "filesystem"
	S3        = "s3"
)

type IStorage interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	GetName() string
	ServeFile(w http.ResponseWriter, r *http.Request, filePath string)
}

type Storage struct {
	storage IStorage
	config  *config.Config
}

func New(c *config.Config) (*Storage, error) {
	var err error
	s := Storage{config: c}
	switch c.Storage {
	case FSStorage:
		s.storage, err = fs.New(c)
	case S3:
		s.storage, err = s3.New(c)
	}
	if s.storage == nil {
		log.Println("couldn't connect to storage")
		return nil, err
	}
	return &s, nil
}

func (s *Storage) Handler() *Storage {
	return s
}

func (s *Storage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.storage.ServeHTTP(w, r)
}
