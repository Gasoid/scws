package s3

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"scws/common/config"

	"github.com/araddon/gou"
	"github.com/lytics/cloudstorage"
	"github.com/lytics/cloudstorage/awss3"
)

const (
	healthPath = "/_/health"
)

func New(c *config.Config) (*S3Storage, error) {
	s3Config := config.S3Config{}
	s3Config.ParseEnv()
	s := S3Storage{
		config: &s3Config,
	}

	conf := &cloudstorage.Config{
		Type:       awss3.StoreType,
		AuthMethod: awss3.AuthAccessKey,
		Bucket:     s.config.Bucket,
		Settings:   make(gou.JsonHelper),
		Region:     os.Getenv("AWS_REGION"),
	}
	conf.Settings[awss3.ConfKeyAccessKey] = os.Getenv("AWS_ACCESS_KEY_ID")
	conf.Settings[awss3.ConfKeyAccessSecret] = os.Getenv("AWS_SECRET_ACCESS_KEY")
	store, err := cloudstorage.NewStore(conf)
	if err != nil {
		log.Println("s3.New", err.Error())
		return nil, err
	}
	s.store = store
	return &s, nil
}

type S3Storage struct {
	config *config.S3Config
	//scwsConfig *config.Config
	store cloudstorage.Store
	index string
}

type object struct {
	prefix string
	store  cloudstorage.Store
	index  string
}

func (o *object) getObject(name string) (cloudstorage.Object, error) {
	obj, err := o.store.Get(context.Background(), path.Join(o.prefix, name))
	if err != nil {
		log.Println("s3.getObject", err.Error())
		return nil, err
	}
	return obj, nil
}

func (o *object) Open(name string) (http.File, error) {
	obj, err := o.getObject(name)
	if err != nil {
		//log.Println("s3.Open", err.Error())
		obj, err = o.getObject(o.index)
		if err != nil {
			log.Println("s3.Open", err.Error())
			return nil, err
		}
	}
	f, err := obj.Open(cloudstorage.ReadOnly)
	if err != nil {
		log.Println("s3.Open", err.Error())
		return nil, err
	}
	return f, nil
}

func (s *S3Storage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == healthPath {
		healthHandler(w, r)
		return
	}
	o := s.newObject()

	http.FileServer(o).ServeHTTP(w, r)
}

func (s *S3Storage) newObject() *object {
	return &object{
		prefix: s.config.Prefix,
		store:  s.store,
		index:  s.index,
	}
}

func (s *S3Storage) GetName() string {
	return "s3"
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
