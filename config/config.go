package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"scws/vault"
	"strings"
)

const (
	tagName      = "default"
	varNameTempl = "SCWS_%s%s"
)

// Config configures scws
type Config struct {
	Port           string `default:"8080"`
	Domain         string
	Storage        string `default:"filesystem"`
	IndexHtml      string `default:"index.html" name:"index_html"`
	VaultAddress   string `name:"vault_address"`
	VaultToken     string `name:"vault_token"`
	VaultPaths     string `name:"vault_paths"`
	SettingsPrefix string `default:"SCWS_SETTINGS_VAR_" name:"settings_prefix"`
	SettingsPath   string `default:"/_/settings" name:"settings_path"`
}

// FsConfig configures Fs storage
type FsConfig struct {
	Root string `default:"/www/"`
}

// S3Config configures S3 storage
type S3Config struct {
	Bucket             string `default:""`
	Prefix             string `default:"/"`
	AwsAccessKeyID     string `name:"Aws_Access_Key_ID" vault:"enabled" vault_alt_key:"access_key"`
	AwsSecretAccessKey string `name:"Aws_Secret_Access_Key" vault:"enabled" vault_alt_key:"secret_key"`
	AwsRegion          string `name:"Aws_Region"`
}

func New() *Config {
	c := Config{}
	c.ParseEnv()

	if c.IsVaultEnabled() {
		err := vault.Init(c.VaultAddress, c.VaultToken)
		if err != nil {
			log.Println("config.New:", err.Error())
		}
	}
	return &c
}

func getEnvVar(name, prefix string) string {
	if prefix != "" {
		prefix = prefix + "_"
	}
	return os.Getenv(fmt.Sprintf(varNameTempl, prefix, strings.ToUpper(name)))
}

type config interface {
	ParseEnv() error
}

func parseEnv(c config, prefix string) {
	t := reflect.ValueOf(c).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		typeField := t.Type().Field(i)
		envName := typeField.Tag.Get("name")
		if envName == "" {
			envName = typeField.Name
		}
		value := getEnvVar(envName, prefix)
		if value == "" {
			tagValue := typeField.Tag.Get(tagName)
			if tagValue != "" {
				tag := strings.TrimSpace(tagValue)
				value = tag
			}

		}
		f.Set(reflect.ValueOf(value))
	}
}

func (c *FsConfig) ParseEnv() error {
	parseEnv(c, "FS")
	return nil
}

func (c *S3Config) ParseEnv() error {
	parseEnv(c, "S3")
	if c.AwsAccessKeyID == "" || c.AwsSecretAccessKey == "" {
		return errors.New("couldn't read AWS credentials: AwsAccessKeyID or AwsSecretAccessKey")
	}
	return nil
}

func (c *Config) ParseEnv() error {
	parseEnv(c, "")
	return nil
}

func (c *Config) GetAddr() string {
	return fmt.Sprintf(":%s", c.Port)
}

func (c *Config) IsVaultEnabled() bool {
	return c.VaultAddress != "" && c.VaultPaths != "" && c.VaultToken != ""
}

func (c *S3Config) GetVaultSecrets(paths string) error {
	pathList := strings.Split(paths, ",")
	for _, p := range pathList {
		secrets, err := vault.Secrets(p)
		if err != nil {
			return fmt.Errorf("couldn't get secrets from vault server: %v", err)
		}
		setConfigVars(c, "S3", secrets)
	}
	return nil
}

func setConfigVars(c config, prefix string, secrets map[string]string) {
	t := reflect.ValueOf(c).Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		typeField := t.Type().Field(i)
		checkVault := typeField.Tag.Get("vault")
		if checkVault != "enabled" {
			continue
		}
		altKey := typeField.Tag.Get("vault_alt_key")
		key := typeField.Tag.Get("name")
		if key == "" {
			key = typeField.Name
		}
		fullKey := fmt.Sprintf(varNameTempl, prefix, strings.ToUpper(key))
		for _, k := range []string{fullKey, altKey} {
			if v, ok := secrets[k]; ok {
				f.Set(reflect.ValueOf(v))
			}
		}

	}
}
