package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewScwsHandler(t *testing.T) {
	h := dummyHandler("")
	testMux := newScwsHandler(map[string]http.Handler{"/": &h})
	assert.NotNil(t, testMux)
}

func TestScwsHandler(t *testing.T) {
	h := dummyHandler("")
	testMux := ScwsHandler{routes: map[string]http.Handler{"/": &h}}
	testHandler := testMux.Handler(&h)
	assert.NotNil(t, testHandler)
}
