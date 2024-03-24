package controller

import (
	"errors"
	"fmt"
	"log/slog"
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

func (c *Controller) LoadConfig(file string) error {
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	config := Config{}
	if err := yaml.Unmarshal(b, &config); err != nil {
		return err
	}
	switch config.ClientType {
	case "oss":
		config.AccessKey = os.Getenv(environmentKeyOssAk)
		config.SecretKey = os.Getenv(environmentKeyOssSk)
		if config.AccessKey == "" || config.SecretKey == "" {
			err := fmt.Errorf("environment key [%s] or [%s] is not set", environmentKeyOssSk, environmentKeyOssSk)
			slog.Error(err.Error())
			return err
		}
	default:
		err := errors.New("unsupported client type")
		slog.Error(err.Error())
		panic(err.Error())
	}

	c.config = &config
	return nil
}
