package controller_test

import (
	"fmt"
	"os"
	"path/filepath"
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

func createDirectoriesAndFile(filePath, content string) error {
	dir := filepath.Dir(filePath)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directories: %v", err)
	}

	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}

	return nil
}

func TestController_UploadDownload(t *testing.T) {
	type T struct {
		filePath string
		content  string
	}
	const uploadDir = "./uploads"
	const downloadDir = "./downloads"

	cases := []T{
		{
			filePath: "aa/bb/cc.txt",
			content:  "cc",
		},
		{
			filePath: "xx.txt",
			content:  "xx",
		},
		{
			filePath: "aa/ddd.txt",
			content:  "ddd",
		},
		{
			filePath: "aa/eee.txt",
			content:  "eee",
		},
	}

	for _, cc := range cases {
		err := createDirectoriesAndFile(filepath.Join(uploadDir, cc.filePath), cc.content)
		assert.NoError(t, err)

	}

	c := newTestCtrl()
	err := c.Upload(controller.UploadOptions{
		Prefix:     prefix,
		EncryptKey: secretKey,
		Paths:      []string{uploadDir},
	})
	assert.NoError(t, err)

	for _, cc := range cases {
		key := filepath.Join(prefix, uploadDir, cc.filePath)
		err = c.Download(controller.DownloadOptions{
			OutputDir:  downloadDir,
			DecryptKey: secretKey,
			S3keys:     []string{key},
		})
		assert.NoError(t, err)

		outputPath := filepath.Join(downloadDir, key)
		content, err := os.ReadFile(outputPath)
		assert.NoError(t, err)
		assert.Equal(t, cc.content, string(content))
	}

	_ = os.RemoveAll(uploadDir)
	_ = os.RemoveAll(downloadDir)
}
