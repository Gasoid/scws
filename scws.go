package main

import (
	"log"
	"net/http"
	"os"
	"scws/config"
	"scws/settings"
	"scws/storage"
	"scws/tracing"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	metricsPath           = "/_/metrics"
	healthPath            = "/_/health"
	rootPath              = "/"
	ifModifiedSinceHeader = "If-Modified-Since"
)

func Run() {
	c := config.New()
	s, err := storage.New(c)
	if err != nil {
		log.Println("Run", err.Error())
		return
	}
	closer, err := tracing.JaegerInit()
	if err != nil {
		log.Println("Run", err.Error())
		return
	}
	defer closer.Close()
	setts := settings.New(c.SettingsPrefix, os.Environ)
	scwsHandler := newScwsHandler(
		map[string]http.Handler{
			c.SettingsPath: setts.Handler(),
			metricsPath:    promhttp.Handler(),
			healthPath:     s.HealthProbe(),
			rootPath:       s.Handler(),
		}, rootPath)
	srv := newServer(c.GetAddr(), scwsHandler)
	catchSignal(srv, setts)
	log.Printf("Starting server on %s", c.GetAddr())
	log.Fatal(srv.ListenAndServe())
}

func newServer(addr string, handler http.Handler) *http.Server {
	srv := &http.Server{
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      handler,
		Addr:         addr,
	}
	return srv
}
