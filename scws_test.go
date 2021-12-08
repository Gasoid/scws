package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyHandler string

func (d *dummyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "dummy handler works")
}

func (d *dummyHandler) Handler() http.Handler {
	return d
}

func (d *dummyHandler) HealthProbe() http.Handler {
	return d
}

func TestNewScwsMux(t *testing.T) {
	h := dummyHandler("")
	testMux := newScwsMux(&h, &h)
	assert.NotNil(t, testMux)
}

func TestScwsHandler(t *testing.T) {
	h := dummyHandler("")
	testHandler := scwsHandler(&h)
	assert.NotNil(t, testHandler)
}

func TestNewServer(t *testing.T) {
	h := dummyHandler("")
	server := newServer("0.0.0.0:8080", &h)
	assert.NotNil(t, server)
}

// func TestRun(t *testing.T) {
// 	assert.NotPanics(t, func() {
// 		Run()
// 	})
// }
