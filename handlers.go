package main

import (
	"net/http"
)

func newScwsHandler(routes map[string]http.Handler, root string) http.Handler {
	scwsHandler := &ScwsHandler{routes, root}
	scwsMux := http.DefaultServeMux
	for k, v := range routes {
		scwsMux.Handle(k, v)
	}
	return scwsHandler.Handler(scwsMux)
}

type ScwsHandler struct {
	routes   map[string]http.Handler
	rootPath string
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
		if _, ok := scwsHandler.routes[r.URL.Path]; ok {
			writer.Flush()
			return
		}
		if writer.Status() == http.StatusNotFound {
			writer.reset(w)
			writer.status = http.StatusOK
			r.URL.Path = scwsHandler.rootPath
			h.ServeHTTP(writer, r)
		}
		writer.Flush()
	}
	return http.HandlerFunc(fn)
}
