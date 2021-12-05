package settings

import (
	"log"
	"scws/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c := config.New()
	err := c.ParseEnv()
	assert.NoError(t, err)
	assert.NotEmpty(t, c.SettingsPrefix)
	log.Println(c.SettingsPrefix)
	settings := New(c)
	v := settings.vars["VAR1"]
	assert.NotEmpty(t, v)
}
