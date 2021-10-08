package storage

import (
	"log"
	"net/http"
	"scws/common/config"
	"scws/storage/fs"
	"scws/storage/s3"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const (
	FSStorage  = "filesystem"
	S3         = "s3"
	serverName = "static-backend"
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

func (s *Storage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tracer := opentracing.GlobalTracer()
	var span opentracing.Span
	if tracer != nil {
		spanCtx, err := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil {
			span = tracer.StartSpan("storage.ServeHTTP")
		} else {
			span = tracer.StartSpan("storage.ServeHTTP", ext.RPCServerOption(spanCtx))
		}
		defer span.Finish()
		span.SetTag("http.status_code", http.StatusOK)
		span.SetTag("http.url", r.URL.Path)
		span.SetTag("storage", s.storage.GetName())
	}
	if strings.HasSuffix(r.URL.Path, "/") || r.URL.Path == "/" {
		s.storage.ServeFile(w, r, s.config.IndexHtml)
	} else {
		s.storage.ServeHTTP(w, r)
	}
	log.Println(r.URL.Path, r.RemoteAddr, w.Header().Get("status"))
}
