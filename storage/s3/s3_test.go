package s3

import (
	"scws/common/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	//utils.SetupTests()
	c := config.New()
	storage, err := New(c)
	assert.NoError(t, err)
	assert.NotNil(t, storage)
}

func TestNewObject(t *testing.T) {
	//utils.SetupTests()
	c := config.New()
	storage, err := New(c)
	assert.NoError(t, err)
	o := storage.newObject()
	name := "README.md"
	cloudObject, err := o.getObject(name)
	assert.NoError(t, err)
	assert.NotNil(t, cloudObject)
}
