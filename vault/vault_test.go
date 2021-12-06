package vault

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit_Err(t *testing.T) {
	err := Init("false-address", "token")
	assert.Error(t, err)
}
