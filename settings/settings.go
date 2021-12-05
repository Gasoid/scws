package settings

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"scws/config"
	"strings"
)

const (
	varSep = "="
)

func New(c *config.Config) *Settings {
	setts := Settings{
		vars: map[string]string{},
	}
	for _, envVar := range os.Environ() {
		kv := strings.Split(envVar, varSep)
		if strings.HasPrefix(kv[0], c.SettingsPrefix) && len(kv) == 2 {
			key := strings.Replace(kv[0], c.SettingsPrefix, "", 1)
			setts.vars[key] = kv[1]
		}
	}
	return &setts
}

type Settings struct {
	vars map[string]string
}

func (setts *Settings) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	content, _ := json.Marshal(setts.vars)
	fmt.Fprint(w, string(content))
}

func (setts *Settings) Handler() *Settings {
	return setts
}
