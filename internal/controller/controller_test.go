package controller_test

import (
	"testing"

	"github.com/linlanniao/soss/internal/controller"
	"github.com/linlanniao/soss/internal/filehandler"
	"github.com/linlanniao/soss/internal/s3clients/ossclient"
	"github.com/stretchr/testify/assert"
)

const (
	prefix    = "tester/"
	secretKey = "p@ssW0rd"
)

func newTestCtrl() *controller.Controller {
	config := controller.NewConfig("../../config.yaml")
	s3Client := ossclient.NewClient(config.Endpoint, config.AccessKey, config.SecretKey)
	fileHandler := filehandler.NewFileHandler()

	c := controller.NewController(config, s3Client, fileHandler, controller.WithDefaultLogger())
	return c
}

func TestController_List(t *testing.T) {
	c := newTestCtrl()
	err := c.List(controller.ListOptions{})
	assert.NoError(t, err)
}

func TestController_Upload(t *testing.T) {
	c := newTestCtrl()
	err := c.Upload(controller.UploadOptions{
		Prefix:     prefix,
		EncryptKey: secretKey,
		Paths:      []string{"../../README.md"},
	})
	assert.NoError(t, err)
}

func TestController_Download(t *testing.T) {
	c := newTestCtrl()
	err := c.Download(controller.DownloadOptions{
		OutputDir:  "../../tmpdir",
		DecryptKey: secretKey,
		S3keys:     []string{prefix + "README.md"},
	})
	assert.NoError(t, err)
}
