package main

import (
	"log"
	"net/http"
	"scws/common/config"
	"scws/common/tracing"
	"scws/storage"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	c := config.New()
	storage, err := storage.New(c)
	if err != nil {
		log.Println("main", err.Error())
		return
	}
	closer, err := tracing.JaegerInit()
	if err != nil {
		log.Println("main", err.Error())
		return
	}
	defer closer.Close()
	RunServer(c, storage)
}

func RunServer(c *config.Config, s *storage.Storage) {
	http.Handle("/_/metrics", promhttp.Handler())
	http.Handle("/", s)
	log.Printf("Storage type: %s", c.Storage)
	log.Printf("Listening %s", c.GetAddr())
	log.Fatal(http.ListenAndServe(c.GetAddr(), nil))
}
