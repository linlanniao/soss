package controller

import (
	"errors"
	"fmt"
	"os"

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

func NewConfig(configFile string) *Config {
	b, err := os.ReadFile(configFile)
	if err != nil {
		panic(err.Error())
	}
	config := Config{}
	if err := yaml.Unmarshal(b, &config); err != nil {
		panic(err.Error())
	}
	switch config.ClientType {
	case "oss":
		config.AccessKey = os.Getenv(environmentKeyOssAk)
		config.SecretKey = os.Getenv(environmentKeyOssSk)
		if config.AccessKey == "" || config.SecretKey == "" {
			err := fmt.Errorf("environment key [%s] or [%s] is not set", environmentKeyOssSk, environmentKeyOssSk)
			panic(err.Error())
		}
	default:
		err := errors.New("unsupported s3Client type")
		panic(err.Error())
	}

	return &config
}
