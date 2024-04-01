package controller

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/linlanniao/soss/pkg/utils"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AccessKey  string `yaml:"-" json:"-"`
	SecretKey  string `yaml:"-" json:"-"`
	ClientType string `yaml:"client_type" json:"client_type"`
	Endpoint   string `yaml:"endpoint" json:"endpoint"`
	Bucket     string `yaml:"bucket" json:"bucket"`
}

const (
	environmentKeyOssAk = "S3_ACCESS_KEY_ID"
	environmentKeyOssSk = "S3_ACCESS_KEY_SECRET"
)

func NewConfig() (*Config, error) {
	c := &Config{}

	// try to load ak / sk from environment
	c.AccessKey = os.Getenv(environmentKeyOssAk)
	c.SecretKey = os.Getenv(environmentKeyOssSk)

	if c.AccessKey == "" || c.SecretKey == "" {
		err := fmt.Errorf("environment key [%s] or [%s] is not set", environmentKeyOssSk, environmentKeyOssSk)
		return nil, err
	}

	updateConfigFromFile := func(file string, configToUpdate *Config) (*Config, error) {
		b, err := os.ReadFile(file)
		if err != nil {
			return nil, err
		}
		fileCfg := Config{}
		if err := yaml.Unmarshal(b, &fileCfg); err != nil {
			return nil, err
		}

		configToUpdate.Bucket = fileCfg.Bucket
		configToUpdate.Endpoint = fileCfg.Endpoint
		configToUpdate.ClientType = fileCfg.ClientType

		return configToUpdate, nil
	}

	// try to load config from./config.yaml
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if f := filepath.Join(dir, defaultConfigFileName); utils.IsFile(f) {
		if updatedCfg, err := updateConfigFromFile(f, c); err == nil {
			return updatedCfg, nil
		}

	}

	// try to load config from $HOME/.soss/config.yaml
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	if f := filepath.Join(home, ".soss", defaultConfigFileName); utils.IsFile(f) {
		if updatedCfg, err := updateConfigFromFile(f, c); err == nil {
			return updatedCfg, nil
		}
	}

	return c, nil
}

const (
	defaultConfigFileName = "config.yaml"
)

func (c *Config) Validate() error {
	switch c.ClientType {
	case "oss":
		return nil
	default:
		return errors.New("invalid client type")
	}
}
