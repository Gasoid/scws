package settings

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	varSep = "="
)

func New(prefix string, getFuncKeys func() []string) *Settings {
	setts := Settings{
		vars:    map[string]string{},
		prefix:  prefix,
		getKeys: getFuncKeys,
	}
	setts.loadVars()
	return &setts
}

type Settings struct {
	vars    map[string]string
	prefix  string
	getKeys func() []string
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

func (setts *Settings) Reload() {
	setts.vars = map[string]string{}
	setts.loadVars()
}

func (setts *Settings) loadVars() {
	for _, envVar := range setts.getKeys() {
		log.Println(envVar)
		kv := strings.Split(envVar, varSep)
		if strings.HasPrefix(kv[0], setts.prefix) && len(kv) == 2 {
			key := strings.Replace(kv[0], setts.prefix, "", 1)
			setts.vars[key] = kv[1]
		}
	}
}
