package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

const (
	tagName = "default"
)

// Config configures scws
type Config struct {
	Port      string `default:"8080"`
	Domain    string
	Storage   string `default:"filesystem"`
	IndexHtml string `default:"index.html"`
}

// FsConfig configures Fs storage
type FsConfig struct {
	Root string `default:"/www/"`
}

// S3Config configures S3 storage
type S3Config struct {
	Bucket string `default:""`
	Prefix string `default:"/"`
}

func New() *Config {
	c := Config{}
	c.ParseEnv()
	return &c
}

func getEnvVar(name string, prefix string) string {
	if prefix != "" {
		prefix = prefix + "_"
	}
	return os.Getenv(fmt.Sprintf("SCWS_%s%s", prefix, strings.ToUpper(name)))
}

func (c *FsConfig) ParseEnv() error {
	t := reflect.ValueOf(c).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		typeField := t.Type().Field(i)
		value := getEnvVar(typeField.Name, "FS")
		if value == "" {
			tagValue := typeField.Tag.Get(tagName)
			if tagValue != "" {
				tag := strings.TrimSpace(tagValue)
				value = tag
			}

		}
		f.Set(reflect.ValueOf(value))
	}
	return nil
}

func (c *S3Config) ParseEnv() error {
	t := reflect.ValueOf(c).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		typeField := t.Type().Field(i)
		value := getEnvVar(typeField.Name, "S3")
		if value == "" {
			tagValue := typeField.Tag.Get(tagName)
			if tagValue != "" {
				tag := strings.TrimSpace(tagValue)
				value = tag
			}

		}
		f.Set(reflect.ValueOf(value))
	}
	return nil
}

func (c *Config) ParseEnv() error {
	t := reflect.ValueOf(c).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		typeField := t.Type().Field(i)
		value := getEnvVar(typeField.Name, "")
		if value == "" {
			tagValue := typeField.Tag.Get(tagName)
			if tagValue != "" {
				tag := strings.TrimSpace(tagValue)
				value = tag
			}

		}
		f.Set(reflect.ValueOf(value))
	}
	return nil
}

func (c *Config) GetAddr() string {
	return fmt.Sprintf(":%s", c.Port)
}
