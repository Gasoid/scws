package settings

import (
	"log"
	"os"
	"scws/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	c := config.New()
	c.ParseEnv()
	assert.NotEmpty(t, c.SettingsPrefix)
	log.Println(c.SettingsPrefix)
	settings := New(c.SettingsPrefix, os.Environ)
	v := settings.vars["VAR1"]
	assert.NotEmpty(t, v)
}
