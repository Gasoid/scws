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

const (
	metricsPath  = "/_/metrics"
	settingsPath = "/_/settings"
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
	scwsSettings := settings.New(c)
	scwsMux := newScwsMux(s.Handler(), scwsSettings.Handler())
	srv := newServer(c, scwsHandler(scwsMux))
	catchSignal(srv, c)
	log.Printf("Starting server on %s", c.GetAddr())
	log.Fatal(srv.ListenAndServe())
}

func newScwsMux(storageHandler http.Handler, settingsHandler http.Handler) *http.ServeMux {
	scwsMux := http.DefaultServeMux
	scwsMux.Handle(metricsPath, promhttp.Handler())
	scwsMux.Handle(settingsPath, settingsHandler)
	scwsMux.Handle("/", storageHandler)
	return scwsMux
}

func newServer(c *config.Config, handler http.Handler) *http.Server {
	srv := &http.Server{
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      handler,
		Addr:         c.GetAddr(),
	}
	return srv
}

func scwsHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		writer := &responseWriter{
			ResponseWriter: w,
			status:         defaultStatus,
		}
		h.ServeHTTP(writer, r)
		logRequest(writer, r)
		traceRequest(writer, r)
		if r.URL.Path == "/" || r.URL.Path == settingsPath || r.URL.Path == metricsPath {
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
