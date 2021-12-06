package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogRequest(t *testing.T) {
	testWriter := httptest.NewRecorder()
	writer := &responseWriter{
		ResponseWriter: testWriter,
		status:         defaultStatus,
	}
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	assert.NotPanics(t, func() {
		logRequest(writer, r)
	})
}
