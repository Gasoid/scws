package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckFile(t *testing.T) {
	err := checkFile("index.html")
	assert.Error(t, err)
}

func TestHealthProbe(t *testing.T) {
	storage, err := New("index.html")
	assert.NoError(t, err)
	err = storage.HealthProbe()
	assert.Error(t, err)
}

func TestIndexPath(t *testing.T) {
	storage, err := New("index.html")
	assert.NoError(t, err)
	assert.Equal(t, storage.indexPath(), "/www/index.html")
	storage.index = "dir1/index.html"
	assert.Equal(t, storage.indexPath(), "/www/dir1/index.html")
}
