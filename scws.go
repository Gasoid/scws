package main

import (
	"log"
	"net/http"
	"scws/common/config"
	"scws/common/settings"
	"scws/common/tracing"
	"scws/storage"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	srv := newServer(c, s)
	log.Printf("Starting server on %s", c.GetAddr())
	log.Fatal(srv.ListenAndServe())
}

func newServer(c *config.Config, s *storage.Storage) *http.Server {
	scwsMux := http.DefaultServeMux
	scwsMux.Handle("/_/metrics", promhttp.Handler())
	scwsMux.Handle("/_/settings", settings.New(c))
	scwsMux.Handle("/", s)
	var handler http.Handler = scwsMux
	handler = scwsHandler(handler, c)

	srv := &http.Server{
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      handler,
		Addr:         c.GetAddr(),
	}
	return srv
}

func scwsHandler(h http.Handler, c *config.Config) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		writer := &responseWriter{
			ResponseWriter: w,
			status:         defaultStatus,
		}
		h.ServeHTTP(writer, r)
		logRequest(writer, r)
		traceRequest(writer, r)
		if r.URL.Path == "/" || r.URL.Path == "/_/settings" || r.URL.Path == "/_/metrics" {
			writer.Flush()
			return
		}
		if writer.Status() == http.StatusNotFound {
			writer.reset(w)
			writer.status = http.StatusOK
			r.URL.Path = "/"
			h.ServeHTTP(writer, r)
		}
		writer.Flush()
	}
	return http.HandlerFunc(fn)
}
