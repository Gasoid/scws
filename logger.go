package main

import (
	"log"
	"net/http"
	"time"
)

type logEntity struct {
	Request    *http.Request
	Method     string
	StatusCode int
	Path       string
	RemoteAddr string
	Size       int
	Duration   time.Duration
	UserAgent  string
}

func logRequest(w ResponseWriter, r *http.Request) {
	l := &logEntity{
		Request:    r,
		Method:     r.Method,
		Path:       r.URL.Path,
		UserAgent:  r.Header.Get("User-Agent"),
		StatusCode: w.Status(),
		RemoteAddr: r.RemoteAddr,
		Size:       w.Size(),
	}
	l.print()
}

func (l *logEntity) print() {
	log.Printf("%s %s %s %d %d %s", l.RemoteAddr, l.Method, l.Path, l.StatusCode, l.Size, l.UserAgent)
}
