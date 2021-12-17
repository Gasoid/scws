package main

import (
	"net/http"
)

func newScwsHandler(routes map[string]http.Handler) http.Handler {
	scwsHandler := &ScwsHandler{routes: routes}
	scwsMux := http.DefaultServeMux
	for k, v := range routes {
		scwsMux.Handle(k, v)
	}
	scwsHandler.metrics = metrics()
	return scwsHandler.Handler(scwsMux)
}

type ScwsHandler struct {
	routes  map[string]http.Handler
	metrics func(w ResponseWriter, r *http.Request)
}

func (scwsHandler *ScwsHandler) Handler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		writer := &responseWriter{
			ResponseWriter: w,
			status:         defaultStatus,
		}
		if r.Header.Get(ifModifiedSinceHeader) != "" && r.Method == http.MethodGet {
			r.Header.Del(ifModifiedSinceHeader)
		}
		h.ServeHTTP(writer, r)
		logRequest(writer, r)
		traceRequest(writer, r)
		scwsHandler.metrics(writer, r)
		if _, ok := scwsHandler.routes[r.URL.Path]; ok {
			writer.Flush()
			return
		}
		writer.Flush()
	}
	return http.HandlerFunc(fn)
}
