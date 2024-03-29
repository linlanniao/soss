package controller

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

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
	environmentKeyOssAk = "OSS_ACCESS_KEY_ID"
	environmentKeyOssSk = "OSS_ACCESS_KEY_SECRET"
)

// NewConfigFromFile creates a new Config object from the given YAML file.
// The file should contain the necessary configuration for the S3 client,
// including the access key ID and secret access key, as well as the S3 bucket
// name and endpoint.
func NewConfigFromFile(configFile string) (*Config, error) {
	b, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	config := Config{}
	if err := yaml.Unmarshal(b, &config); err != nil {
		return nil, err
	}
	switch config.ClientType {
	case "oss":
		config.AccessKey = os.Getenv(environmentKeyOssAk)
		config.SecretKey = os.Getenv(environmentKeyOssSk)
		if config.AccessKey == "" || config.SecretKey == "" {
			err := fmt.Errorf("environment key [%s] or [%s] is not set", environmentKeyOssSk, environmentKeyOssSk)
			return nil, err
		}
	default:
		err := errors.New("unsupported s3Client type")
		return nil, err
	}

	return &config, nil
}

const (
	defaultConfigFileName = "config.yaml"
)

func isFile(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// TryGetConfig attempts to load the S3 configuration from the default
// locations, first from./config.yaml and then from $HOME/.soss/config.yaml.
// If no configuration file is found, an error is returned.
func TryGetConfig() (*Config, error) {
	// 1. try to load config from./config.yaml
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if f := filepath.Join(dir, defaultConfigFileName); isFile(f) {
		return NewConfigFromFile(f)
	}

	// 2. try to load config from $HOME/.soss/config.yaml
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	if f := filepath.Join(home, ".soss", defaultConfigFileName); isFile(f) {
		return NewConfigFromFile(f)
	}

	// 3. return an error
	return nil, errors.New("no config file found, please check ./config.yaml or $HOME/.soss/config.yaml")
}
