package s3

import (
	"scws/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	s3 = "s3"
)

func TestNew(t *testing.T) {
	c := config.New()
	if c.Storage != s3 {
		t.Skip("no aws creds, skipping tests")
	}
	storage, err := New(c.IsVaultEnabled(), c.VaultPaths)
	assert.NoError(t, err)
	assert.NotNil(t, storage)
}

func TestNewObject(t *testing.T) {
	c := config.New()
	if c.Storage != s3 {
		t.Skip("no aws creds, skipping tests")
	}
	storage, err := New(c.IsVaultEnabled(), c.VaultPaths)
	assert.NoError(t, err)
	o := storage.newObject()
	name := "README.md"
	assert.NotNil(t, o.store)
	assert.NotEmpty(t, o.prefix)
	cloudObject, err := o.getObject(name)
	assert.NoError(t, err)
	assert.NotNil(t, cloudObject)
}
